package weaviate

import (
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
)

const Class_Name = "LinkInfo"

func createLinkInfoClass() error {
	dataType := []string{schema.DataTypeText.String()}

	emptyClass := &models.Class{
		Class:       Class_Name,
		Description: "系统链接地址",
		Properties: []*models.Property{
			{DataType: dataType, Name: "title"},
			{DataType: dataType, Name: "description"},
			{DataType: dataType, Name: "url"},
		},
		VectorConfig: map[string]models.VectorConfig{
			"description": {
				Vectorizer:      getVectorizer("description"),
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

func getVectorizer(proNames ...string) map[string]interface{} {
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
func CreateLinkInfoClass() error {
	do, err := WClient.Schema().ClassExistenceChecker().WithClassName(Class_Name).Do(WContext)
	if err != nil {
		return err
	}
	if do {
		return nil
	}
	return createLinkInfoClass()
}
