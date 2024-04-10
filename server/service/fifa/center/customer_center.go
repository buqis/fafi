package center

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/service/fifa/center/client"
	"sync"
)

type CustomerCenter struct {
	ctx     context.Context
	rw      sync.RWMutex
	list    []client.BaseClient
	once    sync.Once
	channel map[string]chan interface{}
}
