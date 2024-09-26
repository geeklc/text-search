package weaviate

import (
	"context"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/auth"
	"go.uber.org/zap"
	"text-search/src/logger"
)

var WClient *weaviate.Client

var WContext = context.Background()

func InitClient() {
	cfg := weaviate.Config{
		Host:       "192.168.0.184:18080",
		Scheme:     "http",
		AuthConfig: auth.ApiKey{Value: "WVF5YThaHlkYwhGUSmCRgsX3tD5ngdN8pkih"},
		Headers:    nil,
	}
	client, err := weaviate.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	schema, err := client.Misc().MetaGetter().Do(WContext)
	if err != nil {
		panic(err)
	}
	logger.Logger.Info("weaviate init..", zap.Any("mataData", schema))
	WClient = client
}
