package weaviate

import (
	"errors"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"go.uber.org/zap"
	"text-search/src/logger"
)

const Class_Name2 = "OLLAMA_TEST"

func createOllamaTestClass() error {
	dataType := []string{schema.DataTypeText.String()}

	emptyClass := &models.Class{
		Class:       Class_Name2,
		Description: "系统链接地址",
		Properties: []*models.Property{
			{DataType: dataType, Name: "title"},
			{DataType: dataType, Name: "description"},
			{DataType: dataType, Name: "url"},
		},
		VectorConfig: map[string]models.VectorConfig{
			"description": {
				Vectorizer:      getVectorizer1("description"),
				VectorIndexType: "flat",
			},
		},
	}
	err := WClient.Schema().ClassCreator().
		WithClass(emptyClass).
		Do(WContext)
	if err != nil {
		return err
	}
	return nil
}

func getVectorizer1(proNames ...string) map[string]interface{} {
	vectorizer := map[string]interface{}{
		"text2vec-ollama": map[string]interface{}{
			"properties":         proNames,
			"apiEndpoint":        "http://192.168.0.117:11434",
			"model":              "nomic-embed-text",
			"vectorizeClassName": false,
		},
	}
	return vectorizer
}

// 创建数据集合
func CreateOllamaTestClass() error {
	do, err := WClient.Schema().ClassExistenceChecker().WithClassName(Class_Name2).Do(WContext)
	if err != nil {
		return err
	}
	if do {
		return nil
	}
	return createOllamaTestClass()
}

// 保存数据
func SaveOllamaTest(title, description, url string) error {
	do, err := WClient.Data().Creator().WithClassName(Class_Name2).
		WithProperties(map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
		}).Do(WContext)
	if err != nil {
		return err
	}
	logger.Logger.Info("保存数据", zap.Any("数据", do.Object))
	return nil
}

// 使用向量搜索数据
func SearchContentFromOllamaTest(content string) (map[string]models.JSONObject, error) {
	do, err := WClient.GraphQL().Get().
		WithClassName(Class_Name2).
		WithFields(
			graphql.Field{Name: "title"},
			graphql.Field{Name: "description"},
			graphql.Field{Name: "url"},
			graphql.Field{
				Name: "_additional",
				Fields: []graphql.Field{
					{Name: "distance"},
				},
			},
		).
		WithNearText(WClient.GraphQL().NearTextArgBuilder().WithConcepts([]string{content})).
		WithLimit(10).
		Do(WContext)
	if err != nil {
		return nil, err
	}
	if len(do.Errors) > 0 {
		return nil, errors.New(do.Errors[0].Message)
	}
	return do.Data, nil
}
