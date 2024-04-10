package fafi

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service"
)

type ApiGroup struct {
	CenterAPI
}

var (
	centerService = &service.ServiceGroupApp.FIFAServiceGroup.CenterService
)
