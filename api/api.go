package api

import (
	"errors"

	"github.com/intchensc/qzone/api/common"
	"github.com/intchensc/qzone/api/friend"
	"github.com/intchensc/qzone/api/group"
	"github.com/intchensc/qzone/api/message"
	"github.com/intchensc/qzone/api/shuoshuo"
	"github.com/intchensc/qzone/auth"
)

// API 接口聚合器
type API struct {
	*common.BaseAPI
}

// New 创建新的API实例
func New() *API {
	return &API{BaseAPI: &common.BaseAPI{Err: errors.New("未登录无法调用 API")}}
}

// SetLogin 设置登录信息
func (a *API) SetLogin(auth auth.BaseAuth) {
	a.BaseAPI.Unpack(auth.GetCookie())
}

// Friend 获取好友API实例
func (a *API) Friend() *friend.API {
	return &friend.API{BaseAPI: a.BaseAPI}
}

// Group 获取群组API实例
func (a *API) Group() *group.API {
	return &group.API{BaseAPI: a.BaseAPI}
}

// Message 获取历史记录API实例
func (a *API) Message() *message.API {
	return &message.API{BaseAPI: a.BaseAPI}
}

// ShuoShuo 获取说说API实例
func (a *API) ShuoShuo() *shuoshuo.API {
	return &shuoshuo.API{BaseAPI: a.BaseAPI}
}
