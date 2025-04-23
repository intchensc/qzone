package history

import "time"

// QZoneHistoryReq QQ空间历史消息请求结构体
type QZoneHistoryReq struct {
	Uin                int64  `json:"uin"`    // QQ号
	Offset             int64  `json:"offset"` // 偏移量
	Count              int64  `json:"count"`  // 请求数目
	BeginTime          string `json:"begin_time"`
	EndTime            string `json:"end_time"`
	Getappnotification string `json:"getappnotification"`
	Getnotifi          string `json:"getnotifi"`
	HasGetKey          string `json:"has_get_key"`
	Useutf8            string `json:"useutf8"`
	Outputhtmlfeed     string `json:"outputhtmlfeed"`
	Scope              string `json:"scope"`
	Set                string `json:"set"`
	Format             string `json:"format"`
	Gtk                string `json:"g_tk"`
}

// QZoneHistoryItem QQ空间历史消息返回结构体
type QZoneHistoryItem struct {
	SenderQQ        string    // 发送方QQ
	ActionType      string    // 互动类型
	ShuoshuoID      string    // 说说ID
	ShuoshuoContent string    // 说说内容
	Content         string    // 互动内容
	CreateTime      time.Time // 发送的时间
	ImgUrls         []string  // 互动内容的图片
	ShuoshuoImgUrls []string  // 说说内容
	// QZoneImages	[]string // TODO: 可考虑加入表情
}
