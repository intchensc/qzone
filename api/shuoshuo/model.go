package shuoshuo

import "time"

// ShuoShuoPublishResp 发布说说响应结构体
type ShuoShuoPublishResp struct {
	Tid      string // 说说Id
	Code     int64  // 响应状态码，0成功
	Now      int64  // 发布时间戳
	FeedInfo string // 说说页面html元素
	Message  string // ？错误后返回的消息
}

// ShuoShuoResp 说说响应结构体
type ShuoShuoResp struct {
	Uin         int64  // 用户QQ
	Name        string // 用户昵称
	Tid         string // 说说Id
	Content     string // 说说内容
	CreateTime  string // 说说创建时间
	CreatedTime int64  // 说说创建时间戳
	PicTotal    int64  // 图片总数
	Cmtnum      int64  // 评论数量
	Secret      int64  // 是否为私密动态
	Pic         []PicResp
}

// PicResp 说说响应结构体中的图片数据
type PicResp struct {
	PicId      string // 图片Id
	Url1       string // 原图更小
	Url2       string // 原图大小
	Url3       string // 原图指定hw
	Smallurl   string // 缩略图
	Curlikekey string // 链接
	Unilikekey string // 链接
}

// Comment 评论简单结构体，目前支持一级评论
type Comment struct {
	ShuoShuoID string    //当前评论所属的说说ID
	OwnerName  string    //当前评论人的昵称
	OwnerUin   int64     //当前评论人的QQ
	Content    string    //评论内容，为空则是图片评论
	PicContent []string  //图片评论链接
	CreateTime time.Time //发布评论的时间戳
}

// LikeResp 点赞响应结构体
type LikeResp struct {
	Ret int64
	Msg string
}

// UploadImageResp 上传图片响应结构体
type UploadImageResp struct {
	Pre        string // 低分辨率url
	URL        string // 完整url
	Width      int64  // 宽
	Height     int64  // 高
	OriginURL  string // 图片的原始url
	Contentlen int64  // 图片大小（字节）
	Albumid    string
	Lloc       string
	Sloc       string
	Type       int64
	Ret        int64
}

// EmotionPublishRequest 发说说请求体
type EmotionPublishRequest struct {
	CodeVersion    string `json:"code_version"`
	Con            string `json:"con"`
	Feedversion    string `json:"feedversion"`
	Format         string `json:"format"`
	Hostuin        int64  `json:"hostuin"`
	Paramstr       string `json:"paramstr"`
	PicBo          string `json:"pic_bo"`
	PicTemplate    string `json:"pic_template"`
	Qzreferrer     string `json:"qzreferrer"`
	Richtype       string `json:"richtype"`
	Richval        string `json:"richval"`
	SpecialURL     string `json:"special_url"`
	Subrichtype    string `json:"subrichtype"`
	SynTweetVerson string `json:"syn_tweet_verson"`
	ToSign         string `json:"to_sign"`
	UgcRight       string `json:"ugc_right"`
	Ver            string `json:"ver"`
	Who            string `json:"who"`
}

// EmotionPublishVo 发说说响应体
type EmotionPublishVo struct {
	Activity     []interface{} `json:"activity"`
	Attach       interface{}   `json:"attach"`
	AuthFlag     int           `json:"auth_flag"`
	Code         int           `json:"code"`
	Conlist      []Conlist     `json:"conlist"`
	Content      string        `json:"content"`
	Message      string        `json:"message"`
	OurlInfo     interface{}   `json:"ourl_info"`
	PicTemplate  string        `json:"pic_template"`
	Right        int           `json:"right"`
	Secret       int           `json:"secret"`
	Signin       int           `json:"signin"`
	Smoothpolicy Smoothpolicy  `json:"smoothpolicy"`
	Subcode      int           `json:"subcode"`
	T1Icon       int           `json:"t1_icon"`
	T1Name       string        `json:"t1_name"`
	T1Ntime      int           `json:"t1_ntime"`
	T1Source     int           `json:"t1_source"`
	T1SourceName string        `json:"t1_source_name"`
	T1SourceURL  string        `json:"t1_source_url"`
	T1Tid        string        `json:"t1_tid"`
	T1Time       string        `json:"t1_time"`
	T1Uin        int           `json:"t1_uin"`
	ToTweet      int           `json:"to_tweet"`
	UgcRight     int           `json:"ugc_right"`
}

// Conlist 说说文字消息
type Conlist struct {
	Con  string `json:"con"`
	Type int    `json:"type"`
}

// UploadImageRequest 上传图片请求体
type UploadImageRequest struct {
	Albumtype        string `json:"albumtype"`
	BackUrls         string `json:"backUrls"`
	Base64           string `json:"base64"`
	Charset          string `json:"charset"`
	Exttype          string `json:"exttype"`
	Filename         string `json:"filename"`
	HdHeight         string `json:"hd_height"`
	HdQuality        string `json:"hd_quality"`
	HdWidth          string `json:"hd_width"`
	JsonhtmlCallback string `json:"jsonhtml_callback"`
	OutputCharset    string `json:"output_charset"`
	OutputType       string `json:"output_type"`
	PSkey            string `json:"p_skey"`
	PUin             int64  `json:"p_uin"`
	Picfile          string `json:"picfile"`
	Qzonetoken       string `json:"qzonetoken"`
	Qzreferrer       string `json:"qzreferrer"`
	Refer            string `json:"refer"`
	Skey             string `json:"skey"`
	Uin              int64  `json:"uin"`
	UploadHd         string `json:"upload_hd"`
	Uploadtype       string `json:"uploadtype"`
	URL              string `json:"url"`
	Zzpanelkey       string `json:"zzpanelkey"`
	Zzpaneluin       int64  `json:"zzpaneluin"`
}

type Smoothpolicy struct {
	ComswDisableSosoSearch  int `json:"comsw.disable_soso_search"`
	L1SwReadFirstCacheOnly  int `json:"l1sw.read_first_cache_only"`
	L2SwDontGetReplyCmt     int `json:"l2sw.dont_get_reply_cmt"`
	L2SwMixsvrFrdnumPerTime int `json:"l2sw.mixsvr_frdnum_per_time"`
	L3SwHideReplyCmt        int `json:"l3sw.hide_reply_cmt"`
	L4SwReadTdbOnly         int `json:"l4sw.read_tdb_only"`
	L5SwReadCacheOnly       int `json:"l5sw.read_cache_only"`
}

// MsgListRequest 说说列表请求体
type MsgListRequest struct {
	Callback           string `json:"callback"`
	CodeVersion        string `json:"code_version"`
	Format             string `json:"format"`
	Ftype              string `json:"ftype"`
	GTk                string `json:"g_tk"`
	NeedPrivateComment string `json:"need_private_comment"`
	Num                string `json:"num"`
	Pos                string `json:"pos"`
	Replynum           string `json:"replynum"`
	Sort               string `json:"sort"`
	Uin                int64  `json:"uin"`
}

// LikeRequest 空间点赞请求体
type LikeRequest struct {
	Curkey     string `json:"curkey"`
	Opuin      int64  `json:"opuin"`
	Qzreferrer string `json:"qzreferrer"`
	Unikey     string `json:"unikey"`
	Fid        string `json:"fid"`
	From       string `json:"from"`
	Typeid     string `json:"typeid"`
	Appid      string `json:"appid"`
}
