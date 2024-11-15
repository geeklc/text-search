package main

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"text-search/src/logger"
	"text-search/src/weaviate"
	"time"
)

const OLLAMA_HOSTNAME = "http://192.168.0.117:11434"

// 代理处理逻辑
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//如果不是回话请求则直接代理信息
	if path != "/api/chat" {
		proxyRequest(w, r, r.Body, false)
		return
	}
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		logger.Logger.Error("读取请求体报错", zap.Any("", err))
		return
	}
	defer r.Body.Close()
	//把请求数据转换为想要的格式
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		logger.Logger.Error("请求数据转换为想要的格式报错", zap.Any("", err))
		return
	}
	msgs, err := dealRequetMsgs(res)
	if err != nil {
		http.Error(w, "Unable to read body", http.StatusBadRequest)
		logger.Logger.Error("处理本地上下文失败", zap.Any("", err))
		return
	}
	dataByte, err := json.Marshal(msgs)
	if err != nil {
		panic(err)
	}
	proxyRequest(w, r, bytes.NewReader(dataByte), true)
}

// 处理请求信息，预置上下文
func dealRequetMsgs(reqBody map[string]interface{}) (map[string]interface{}, error) {
	messages := reqBody["messages"]
	msgArr := messages.([]interface{})
	var newMsgs []map[string]interface{}
	var lastMsg map[string]interface{}
	for i, msg := range msgArr {
		msgMap := msg.(map[string]interface{})
		if i+1 == len(msgArr) {
			lastMsg = msgMap
			continue
		}
		newMsgs = append(newMsgs, msgMap)
	}
	//获取本地向量库的上下文
	content := lastMsg["content"]
	LocalContexts, err := GetLocalContext(content.(string))
	if err != nil {
		return nil, err
	}
	for _, info := range LocalContexts {
		newMsgs = append(newMsgs, info)
	}
	newMsgs = append(newMsgs, lastMsg)
	reqBody["messages"] = newMsgs
	return reqBody, nil
}

// 获取上下文信息
func GetLocalContext(content string) ([]map[string]interface{}, error) {
	localContext, err := weaviate.SearchContentFromBegTest(content)
	if err != nil {
		return nil, err
	}
	logger.Logger.Debug("查询的向量数据", zap.Any("向量数据", localContext))
	getData := localContext["Get"]
	getDataMap := getData.(map[string]interface{})
	dataArr := getDataMap[weaviate.Class_Name_BEG]
	dataArrMaps := dataArr.([]interface{})
	var returnArr []map[string]interface{}
	for _, data := range dataArrMaps {
		dataMap := data.(map[string]interface{})
		additional := dataMap["_additional"].(map[string]interface{})
		distance := additional["distance"].(float64)
		if distance > 0.25 {
			continue
		}
		description := dataMap["description"].(string)
		url := dataMap["url"].(string)
		info := map[string]interface{}{
			"role":    "assistant",
			"content": "你可以整合如下相关内容回答：" + description + "；相关链接为：" + url,
		}
		returnArr = append(returnArr, info)
	}
	return returnArr, nil
}

// 请求转发
func proxyRequest(w http.ResponseWriter, r *http.Request, reqBody io.Reader, streamFlag bool) {

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(r.Method, OLLAMA_HOSTNAME+r.URL.Path, reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	if streamFlag {
		// 设置响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 流式响应
		buffer := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				// 写入流响应
				w.Write(buffer[:n])
				w.(http.Flusher).Flush() // 刷新响应缓冲区
			}
			if err != nil {
				break
			}
		}
	} else {
		body, _ := io.ReadAll(resp.Body)
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}

}
