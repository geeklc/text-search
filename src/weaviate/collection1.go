package weaviate

import (
	"errors"
	"fmt"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"go.uber.org/zap"
	"text-search/src/logger"
	"text-search/src/req"
)

const Class_Name1 = "HELLO_INFO"

func createLinkInfosClass() error {
	dataType := []string{schema.DataTypeText.String()}

	emptyClass := &models.Class{
		Class:       Class_Name1,
		Description: "系统链接地址",
		Properties: []*models.Property{
			{DataType: dataType, Name: "title"},
			{DataType: dataType, Name: "description"},
			{DataType: dataType, Name: "url"},
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

// 创建数据集合
func CreateLinkInfosClass() error {
	do, err := WClient.Schema().ClassExistenceChecker().WithClassName(Class_Name1).Do(WContext)
	if err != nil {
		return err
	}
	if do {
		return nil
	}
	return createLinkInfosClass()
}

// 保存数据
func SaveLinkInfos(title, description, url string) error {
	embedding, err := GetVectorFromEmbedding(description)
	if err != nil {
		return err
	}
	do, err := WClient.Data().Creator().WithClassName(Class_Name1).
		WithProperties(map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
		}).WithVector(embedding).Do(WContext)
	if err != nil {
		return err
	}
	logger.Logger.Info("保存数据", zap.Any("数据", do.Object))
	return nil
}

// 使用向量搜索数据
func SearchContentVector(content string) (map[string]models.JSONObject, error) {
	embedding, err := GetVectorFromEmbedding(content)
	if err != nil {
		return nil, err
	}
	do, err := WClient.GraphQL().Get().
		WithClassName(Class_Name1).
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
		WithNearVector(WClient.GraphQL().NearVectorArgBuilder().WithVector(embedding)).
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

func GetVectorFromEmbedding(content string) ([]float32, error) {
	data := map[string]string{"model": "nomic-embed-text", "input": content}
	resp, err := req.PostFor(data, "/api/embed")
	if err != nil {
		return nil, err
	}
	i := resp["embeddings"]
	return interfaceToFloat32Slice(i)
}

func interfaceToFloat32Slice(input interface{}) ([]float32, error) {
	// 尝试将输入转换为切片
	if slice, ok := input.([]interface{}); ok {
		var fs []float32
		len1 := 0
		for _, arr := range slice {
			if arr1, ok1 := arr.([]interface{}); ok1 {
				for _, v := range arr1 {
					f := v.(float64)
					// 类型断言为 float32
					fs = append(fs, float32(f))
				}
				len1 += len(arr1)
			}
		}
		return fs, nil
	}
	return nil, fmt.Errorf("输入不是 []interface{} 类型")
}
