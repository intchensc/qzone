package group

// QQGroupReq 获取QQ群请求结构体
type QQGroupReq struct {
	Uin     int64  `json:"uin"`
	Do      string `json:"do"`
	Rd      string `json:"rd"`
	Fupdate string `json:"fupdate"`
	Clean   string `json:"clean"`
	GTk     string `json:"g_tk"`
}

// QQGroupResp 获取QQ群响应结构体
type QQGroupResp struct {
	GroupCode   int64  `json:"groupcode"`    //群号
	GroupName   string `json:"groupname"`    //群名
	TotalMember int64  `json:"total_member"` //群人数
	NotFriends  int64  `json:"notfriends"`   //群里非好友人数
}

// QQGroupMemberReq QQ群非好友请求结构体
type QQGroupMemberReq struct {
	Uin     int64  `json:"uin"` //QQ
	Gid     int64  `json:"gid"` //群号
	Fupdate string `json:"fupdate"`
	Type    string `json:"type"`
	GTk     string `json:"g_tk"`
}

// QQGroupMemberResp QQ群非好友响应结构体
type QQGroupMemberResp struct {
	Uin       int64  `json:"fuin"` //QQ
	NickName  string `json:"name"` //昵称
	AvatarURL string `json:"img"`  //头像
	GroupCode int64  `json:"gid"`  //所属群
}
