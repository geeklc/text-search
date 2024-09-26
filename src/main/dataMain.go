package main

import (
	"go.uber.org/zap"
	"text-search/src/common"
	"text-search/src/logger"
	"text-search/src/weaviate"
)

func main() {
	err := initConf()
	if err != nil {
		logger.Logger.Error("初始化报错", zap.Any("错误信息", err))
	}
	weaviate.InitClient()
	/*//创建class
	err = weaviate.CreateBegTestClass()
	if err != nil {
		logger.Logger.Error("创建class报错", zap.Any("错误信息", err))
		return
	}

	//保存数据
	err = insertData()
	if err != nil {
		logger.Logger.Error("保存数据报错", zap.Any("错误信息", err))
		return
	}*/

	content, err := weaviate.SearchContentFromBegTest("找一个支持weaviate的向量库管理端")
	if err != nil {
		logger.Logger.Error("保存数据报错", zap.Any("错误信息", err))
		return
	}
	logger.Logger.Info("查询结果：", zap.Any("", content))

}

func insertData() error {
	infos := []map[string]string{
		{"title": "qwen2.5大模型地址", "description": "ollama的qwen2.5大模型地址", "url": "https://ollama.com/library/qwen2.5"},
		{"title": "dify环境变量说明", "description": "自建的dify环境变量说明", "url": "https://docs.dify.ai/v/zh-hans/getting-started/install-self-hosted/environments"},
		{"title": "RouteLLM：经济高效的 LLM 路由开源框架", "description": "基于python编写的RouteLLM：经济高效的 LLM 路由开源框架", "url": "https://lmsys.org/blog/2024-07-01-routellm/"},
		{"title": "weaviate的golang客户端", "description": "在github上的weaviate的golang客户端", "url": "https://github.com/weaviate/weaviate-go-client"},
		{"title": "使用weaviate实现向量存储", "description": "使用weaviate实现向量存储", "url": "https://blog.csdn.net/make_progress/article/details/139048756"},
		{"title": "weaviate的环境变量配置", "description": "weaviate的环境变量参数配置", "url": "https://weaviate.io/developers/weaviate/config-refs/env-vars"},
		{"title": "向量库管理端vector-admin", "description": "向量库管理端vector-admin，支持weaviate、Chroma等向量库", "url": "https://github.com/Mintplex-Labs/vector-admin/blob/master/docker/DOCKER.md"},
	}
	for _, info := range infos {
		err := weaviate.SaveBegTest(info["title"], info["description"], info["url"])
		if err != nil {
			return err
		}
	}
	return nil
}

// 项目基础配置初始化
func initConf() error {
	//1、初始化配置文件
	err := common.InitConfig()
	if err != nil {
		return err
	}
	//2、 初始化log配置
	err = logger.Init()
	if err != nil {
		return err
	}
	return nil
}
