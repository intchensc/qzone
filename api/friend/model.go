package friend

// FriendInfoEasyResp 好友简略信息响应结构体
type FriendInfoEasyResp struct {
	Uin       int64  // QQ号
	Groupid   int64  // 分组ID
	GroupName string // 分组名称
	Name      string // 名称
	Remark    string // 备注
	Image     string // 头像
	Online    int64  // 在线状态
}

// FriendInfoDetailResp 好友详细信息响应结构体
type FriendInfoDetailResp struct {
	Uin           int64  `json:"uin"`           // QQ号
	Nickname      string `json:"nickname"`      // 昵称
	Signature     string `json:"signature"`     // 签名
	Avatar        string `json:"avatar"`        // 上古头像
	Sex           int64  `json:"sex"`           // 性别，1男
	Age           int64  `json:"age"`           // 年龄
	Birthyear     int64  `json:"birthyear"`     // 生日年份
	Birthday      string `json:"birthday"`      // 生日月-天
	Country       string `json:"country"`       // 国家
	Province      string `json:"province"`      // 省份
	City          string `json:"city"`          // 城市
	Career        string `json:"career"`        // 职业
	Company       string `json:"company"`       // 公司
	Mailname      string `json:"mailname"`      // 邮件名称
	Mailcellphone string `json:"mailcellphone"` // 邮件绑定手机号
	Mailaddr      string `json:"mailaddr"`      // 邮件地址
}
