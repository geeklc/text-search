package weaviate

import (
	"errors"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
	"go.uber.org/zap"
	"text-search/src/logger"
)

// 保存数据
func SaveLinkInfo(title, description, url string) error {
	do, err := WClient.Data().Creator().WithClassName(Class_Name).
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

// 保存数据
func SearchContent(content string) (map[string]models.JSONObject, error) {
	do, err := WClient.GraphQL().Get().
		WithClassName(Class_Name).
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
		WithLimit(2).
		Do(WContext)
	if err != nil {
		return nil, err
	}
	if len(do.Errors) > 0 {
		return nil, errors.New(do.Errors[0].Message)
	}
	return do.Data, nil
}

func CountData() (map[string]models.JSONObject, error) {
	response, err := WClient.GraphQL().Aggregate().
		WithClassName(Class_Name).
		WithFields(graphql.Field{
			Name: "title",
			Fields: []graphql.Field{
				{Name: "count"},
			},
		}).
		Do(WContext)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}
