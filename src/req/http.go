package req

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"text-search/src/logger"
	"time"
)

// Post请求
func PostFor(data any, uri string) (map[string]interface{}, error) {
	if data == nil {
		logger.Logger.Error("请求数据为空", zap.Any("请求uri：", uri))
		panic(uri + "请求数据为空")
	}
	logger.Logger.Debug("开始http请求.....", zap.Any("请求uri", uri))
	//0、实体转换成json串
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	//1、生成client 参数为默认
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	//2、提交请求
	url := getUrl(uri)
	req, err2 := http.NewRequest(http.MethodPost, url, bytes.NewReader(dataByte))
	//3、添加头信息
	addHeader(req)
	if err2 != nil {
		return nil, err2
	}
	//4、处理返回结果
	resp, err3 := client.Do(req)
	if err3 != nil {
		return nil, err3
	}
	defer resp.Body.Close()
	//5、获取返回的状态码
	status := resp.StatusCode
	if status != 200 {
		logger.Logger.Error("请求响应失败", zap.Any("请求地址：", url),
			zap.Any("返回状态码：", status), zap.Any("请求参数：", data))
		panic("")
	}
	//4、获取返回结果
	body, _ := io.ReadAll(resp.Body)
	resultStr := string(body)
	//5、获取返回格式的json
	var res map[string]interface{}
	err = json.Unmarshal([]byte(resultStr), &res)
	if err != nil {
		logger.Logger.Debug("格式化数据失败：", zap.Any("原数据", resultStr), zap.Any("失败信息", err))
		return nil, err
	}
	return res, nil
}

// 添加头信息
func addHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
}

// 获取请求路径
func getUrl(uri string) string {
	url := "http://192.168.0.117:11434"
	if uri == "" {
		return url
	}
	if !strings.HasSuffix(url, "/") && !strings.HasPrefix(uri, "/") {
		url = url + "/"
	}
	url = url + uri
	return url
}
