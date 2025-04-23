package api

import (
	"errors"

	"github.com/intchensc/qzone/api/common"
	"github.com/intchensc/qzone/api/friend"
	"github.com/intchensc/qzone/api/group"
	"github.com/intchensc/qzone/api/history"
	"github.com/intchensc/qzone/api/shuoshuo"
	"github.com/intchensc/qzone/auth"
)

// API 接口聚合器
type API struct {
	*common.BaseAPI
}

// NewAPI 创建新的API实例
func New() *API {
	return &API{BaseAPI: &common.BaseAPI{Err: errors.New("未登录无法调用 API")}}
}

// SetLogin 设置登录信息
func (a *API) SetLogin(auth auth.BaseAuth) {
	a.BaseAPI.Unpack(auth.GetCookie())
}

// Friend 获取好友API实例
func (a *API) Friend() *friend.FriendAPI {
	return &friend.FriendAPI{BaseAPI: a.BaseAPI}
}

// Group 获取群组API实例
func (a *API) Group() *group.GroupAPI {
	return &group.GroupAPI{BaseAPI: a.BaseAPI}
}

// History 获取历史记录API实例
func (a *API) History() *history.HistoryAPI {
	return &history.HistoryAPI{BaseAPI: a.BaseAPI}
}

// ShuoShuo 获取说说API实例
func (a *API) ShuoShuo() *shuoshuo.ShuoShuoAPI {
	return &shuoshuo.ShuoShuoAPI{BaseAPI: a.BaseAPI}
}
