package weaviate

import (
	"errors"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"go.uber.org/zap"
	"text-search/src/logger"
)

const Class_Name_BEG = "BEG_TEST"

// 直接创建class
func createBegTestClass() error {
	dataType := []string{schema.DataTypeText.String()}

	emptyClass := &models.Class{
		Class:       Class_Name_BEG,
		Description: "系统链接地址",
		Properties: []*models.Property{
			{DataType: dataType, Name: "title"},
			{DataType: dataType, Name: "description"},
			{DataType: dataType, Name: "url"},
		},
		VectorConfig: map[string]models.VectorConfig{
			"description": {
				Vectorizer:      getVectorizer2("description"),
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

// 指定ollama的地址和模型名称
func getVectorizer2(proNames ...string) map[string]interface{} {
	vectorizer := map[string]interface{}{
		"text2vec-ollama": map[string]interface{}{
			"properties":         proNames,
			"apiEndpoint":        "http://192.168.0.117:11434",
			"model":              "znbang/bge:large-zh-v1.5-f32",
			"vectorizeClassName": false,
		},
	}
	return vectorizer
}

// 创建数据集合
func CreateBegTestClass() error {
	do, err := WClient.Schema().ClassExistenceChecker().WithClassName(Class_Name_BEG).Do(WContext)
	if err != nil {
		return err
	}
	if do {
		return nil
	}
	return createBegTestClass()
}

// 保存数据
func SaveBegTest(title, description, url string) error {
	do, err := WClient.Data().Creator().WithClassName(Class_Name_BEG).
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
func SearchContentFromBegTest(content string) (map[string]models.JSONObject, error) {
	do, err := WClient.GraphQL().Get().
		WithClassName(Class_Name_BEG).
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
		WithLimit(1).
		Do(WContext)
	if err != nil {
		return nil, err
	}
	if len(do.Errors) > 0 {
		return nil, errors.New(do.Errors[0].Message)
	}
	return do.Data, nil
}
