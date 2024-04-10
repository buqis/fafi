package fafi

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/fifa/local_data"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/fifa/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CenterAPI struct{}

// 获取账户列表
func (s *CenterAPI) GetAccountList(c *gin.Context) {
	var pageInfo request.PageInfo
	//把页码消息解析成结构体
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//解析jwt中的用户id
	uid := utils.GetUserID(c)
	//根据用户id和页码数获取分页后的账户列表
	list, total, err := centerService.GetAccountInfoList(uid, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败！", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// 账号操作
func (s *CenterAPI) AccountEvent(c *gin.Context) {
	var ae request2.AccountEvent
	err := c.ShouldBindJSON(ae)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(ae, utils.AccountEventVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	uid := utils.GetUserID(c)
	//根据传进来的账户活动，对账户进行启动，删除，停止操作
	err = centerService.AccountEvent(uid, ae)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("操作成功", c)
}

// 添加账号
//
//	func (s *CenterAPI) CreateAccount(c *gin.Context){
//		var input request2.CreateAccount
//
// }
// 根据id查找账户
func (s *CenterAPI) GetAccountById(c *gin.Context) {
	var idInfo request.GetById
	//ShouldBindQuery是用于绑定url地址中的值，等于getQuery
	//当绑定发生错误时（如参数输入错误），不会自动返回400状态码并将Content-Type 被设置为 text/plain; charset=utf-8
	_ = c.ShouldBindQuery(&idInfo)
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	//根据用户id获取账号，返回的是账号读接口和错误信息
	account, err := centerService.GetAccountById(uid, uint(idInfo.ID))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败，原因: "+err.Error(), c)
		return
	}
	if account == nil {
		global.GVA_LOG.Error("获取失败，账号不存在")
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(account, "获取成功", c)
}

// 修改账户信息
func (s *CenterAPI) UpdateAccount(c *gin.Context) {
	var account local_data.Account
	err := c.ShouldBindJSON(&account)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	if err = centerService.UpdateAccount(uid, account); err != nil {
		global.GVA_LOG.Error("修改失败", zap.Error(err))
		response.FailWithMessage("修改失败，原因:"+err.Error(), c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// 设置代理
func (s *CenterAPI) SetProxy(c *gin.Context) {
	var setProxyReq request2.SetProxyReq
	err := c.ShouldBindJSON(&setProxyReq)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	if err = utils.Verify(setProxyReq, utils.SetProxyVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = centerService.SetProxy(uid, setProxyReq); err != nil {
		global.GVA_LOG.Error("设置失败", zap.Error(err))
		response.FailWithMessage("设置失败,原因: "+err.Error(), c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}
