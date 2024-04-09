package client

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
)

// 创建一个客户端
func newBaseClient(ctx context.Context, account local_data.AccountWriterAndReader) *BaseClient {
	baseClient := &BaseClient{
		Account:   account,
		Proxy:     nil,
		Jar:       nil,
		ParentCtx: ctx,
		CliState:  local_data.CliStateOffline,
	}
	return baseClient
}

func NewBaseClient(ctx context.Context, account local_data.AccountWriterAndReader) *BaseClient {
	baseClient := newBaseClient(ctx, account)
	baseClient.Task = baseClient.Start
	return baseClient
}

var User = make(chan struct{}, 50)
