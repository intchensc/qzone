package message

import "time"

// Message QQ空间历史消息返回结构体
type Message struct {
	TID        string      // 说说ID
	OriginTID  string      // 被转发的说说ID
	QQ         string      // 互动人
	Type       MessageType // 互动类型
	Raw        string      // 原始HTMl
	CreateTime time.Time   // 发送的时间
}

type MessageType int

// data-details 是相关说说的地址
// data-feedstype 是动态类型
// data-origtid 如果是转发的，这个是原说说
// data-tid 不是转发，上面的和这个都是本条相关的tid
// data-feedstype="107" 转发艾特人评论 103-转发艾特不评论 。101点赞
const (
	TLike      MessageType = iota // 说说被赞 "赞了我的说说"
	TAt                           // 被艾特 "提到我"
	TCommentAt                    // 评论艾特 “评论提到我”
	TComment                      // 被艾特 "评论"
	TOther                        // 上述以外的
)
