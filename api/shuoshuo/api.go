package shuoshuo

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/intchensc/qzone/api/common"
	"github.com/tidwall/gjson"
)

var (
	// cReLike 点赞响应正则，frameElement.callback();
	cReLike           = regexp.MustCompile(`(?s)frameElement.callback\((.*)\)`)
	cRe               = regexp.MustCompile(`(?s)_Callback\((.*)\)`)
	userQzoneURL      = "https://user.qzone.qq.com"
	msglistURL        = "https://user.qzone.qq.com/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msglist_v6?"
	getCommentsURL    = "https://h5.qzone.qq.com/proxy/domain/taotao.qq.com/cgi-bin/emotion_cgi_msgdetail_v6?uin=%s&pos=%d&num=%d&tid=%s&format=jsonp&g_tk=%s"
	emotionPublishURL = "https://user.qzone.qq.com/proxy/domain/taotao.qzone.qq.com/cgi-bin/emotion_cgi_publish_v6?g_tk=%v"
	uploadImageURL    = "https://up.qzone.qq.com/cgi-bin/upload/cgi_upload_image?g_tk=%v"
	likeURL           = "https://user.qzone.qq.com/proxy/domain/w.qzone.qq.com/cgi-bin/likes/internal_dolike_app?g_tk=%v"
	getLikeListURL    = "https://h5.qzone.qq.com/proxy/domain/users.qzone.qq.com/cgi-bin/likes/get_like_list_app?"
)

type ShuoShuoAPI struct {
	*common.BaseAPI
}

// PublishShuoShuo 发布说说，content文本内容，base64imgList图片数组
func (s *ShuoShuoAPI) Publish(content string, base64imgList []string) (*ShuoShuoPublishResp, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	var (
		picBo       string
		richval     string
		richvalList = make([]string, 0, 9)
		picBoList   = make([]string, 0, 9)
	)

	for _, base64img := range base64imgList {
		uir, err := s.uploadImage(base64img)
		if err != nil {
			log.Println("说说发布失败:", err.Error())
			return nil, err
		}
		picBo, richval, err = s.getPicBoAndRichval(uir)
		if err != nil {
			log.Println("说说发布失败:", err.Error())
			return nil, err
		}
		richvalList = append(richvalList, richval)
		picBoList = append(picBoList, picBo)
	}

	epr := EmotionPublishRequest{
		SynTweetVerson: "1",
		Paramstr:       "1",
		Who:            "1",
		Con:            content,
		Feedversion:    "1",
		Ver:            "1",
		UgcRight:       "1",
		ToSign:         "0",
		Hostuin:        s.Qq,
		CodeVersion:    "1",
		Format:         "json",
		Qzreferrer:     userQzoneURL + "/" + strconv.FormatInt(s.Qq, 10),
	}
	if len(base64imgList) > 0 {
		epr.Richtype = "1"
		epr.Richval = strings.Join(richvalList, "\t")
		epr.Subrichtype = "1"
		epr.PicBo = strings.Join(picBoList, ",")
	}
	url := fmt.Sprintf(emotionPublishURL, s.Gtk2)
	// payload := strings.NewReader(structToStr(epr))
	var data string
	V := strings.NewReader(common.StructToStr(epr))
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		Header("origin", userQzoneURL).
		UserAgent(common.UA).ContentType(common.ContentType).ToString(&data).
		Header("cookie", s.Cookie).
		BodyReader(V).Fetch(context.Background())

	// data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
	// 	WithHeader(map[string]string{
	// 		"referer": userQzoneURL,
	// 		"origin":  userQzoneURL,
	// 		"cookie":  s.Cookie,
	// 	})))
	if err != nil {
		// er := errors.New("说说发布请求错误:" + err.Error())
		log.Println("说说发布失败:", err)
		// return nil, er
	}

	jsonStr := data
	ssp := &ShuoShuoPublishResp{
		Code:    gjson.Get(jsonStr, "code").Int(),
		Tid:     gjson.Get(jsonStr, "tid").String(),
		Now:     gjson.Get(jsonStr, "now").Int(),
		Message: gjson.Get(jsonStr, "message").String(),
	}
	if ssp.Message != "" {
		er := errors.New("说说发布错误:" + ssp.Message)
		log.Println("说说发布失败:", er.Error())
		return nil, er
	}
	return ssp, nil
}

// ShuoShuoList 获取所有说说 实际能访问的说说个数 <= 说说总数(空间仅展示近半年等情况) (有空间访问权限即可)
func (s *ShuoShuoAPI) List(uin int64, num int64, ms int64) (ShuoShuo []*ShuoShuoResp, err error) {
	if s.Err != nil {
		return nil, s.Err
	}
	cnt := num
	t := int(math.Ceil(float64(cnt) / 20.0))
	var i int
	//获取最大数量，控制i的取值
	maxCnt, err := s.Count(uin)
	if err != nil {
		log.Println("说说获取失败:", err.Error())
		return nil, err
	}
	for range t {
		if i >= int(maxCnt) {
			break
		}
		ShuoShuoTemp, err := s.shuoShuoListRaw(uin, 20, i, 0)
		if err != nil {
			log.Println("所有说说获取失败:", err.Error())
			return nil, err
		}
		if len(ShuoShuoTemp) == 0 {
			break
		}
		if len(ShuoShuo) < int(cnt) {
			ShuoShuo = append(ShuoShuo, ShuoShuoTemp[0:min(len(ShuoShuoTemp), int(cnt)-len(ShuoShuo))]...)
			i = i + 20
			time.Sleep(time.Millisecond * time.Duration(ms))
		}
	}
	return ShuoShuo, nil
}

// GetShuoShuoCount 获取用户QQ号为uin的说说总数（有空间访问权限即可）
func (s *ShuoShuoAPI) Count(uin int64) (int64, error) {
	if s.Err != nil {
		return -1, s.Err
	}
	mlr := MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Pos:                "0",
		Num:                "1",
		Replynum:           "0",
		GTk:                s.Gtk2,
		Callback:           "_preloadCallback",
		CodeVersion:        "1",
		Format:             "json",
		NeedPrivateComment: "1",
	}
	url := msglistURL + common.StructToStr(mlr)
	var data string
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		Header("origin", userQzoneURL).
		Header("cookie", s.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"referer": userQzoneURL,
	// 	"origin":  userQzoneURL,
	// 	"cookie":  s.Cookie,
	// })))
	if err != nil {
		// er := errors.New("说说总数请求错误:" + err.Error())
		log.Println("说说总数获取失败:", err)
		// return -1, er
	}
	jsonStr := data
	// 判断是否有访问权限
	forbid := gjson.Get(jsonStr, "message").String()
	if forbid != "" {
		er := errors.New("说说总数响应错误:" + forbid)
		log.Println("说说总数响应失败:", er.Error())
		return -1, er
	}
	cnt := gjson.Get(jsonStr, "total").Int()
	return cnt, nil
}

// GetLevel1CommentCount 获取一级评论总数(限制本人)
func (s *ShuoShuoAPI) GetLevel1CommentCount(tid string) (int64, error) {
	if s.Err != nil {
		return -1, s.Err
	}
	url := fmt.Sprintf(getCommentsURL, strconv.FormatInt(s.Qq, 10), 0, 1, tid, s.Gtk2)
	var data string
	err := requests.
		URL(url).
		Header("cookie", s.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"cookie": s.Cookie,
	// })))
	if err != nil {
		// er := errors.New("说说评论请求错误:" + err.Error())
		log.Println("说说评论请求失败:", err)
		// return -1, er
	}
	r := cRe.FindStringSubmatch(string(data))
	//log.Println("空指针异常测试：" + string(data))
	if len(r) < 2 {
		// er := errors.New("说说评论正则解析错误:" + err.Error())
		log.Println("说说评论请求失败:", err)
		// return -1, er
	}
	jsonRaw := r[1]

	// 说说的一级评论总数
	numOfComments := gjson.Get(jsonRaw, "cmtnum").Int()
	return numOfComments, nil
}

// ShuoShuoCommentList 根据说说ID获取评论（限制本人）
func (s *ShuoShuoAPI) CommentList(tid string, num int64, ms int64) (comments []*Comment, err error) {
	if s.Err != nil {
		return nil, s.Err
	}
	numOfComments := num
	t := int(math.Ceil(float64(numOfComments) / 20.0))
	//获取最大数量，控制i的取值
	maxCnt, err := s.GetLevel1CommentCount(tid)
	if err != nil {
		log.Println("说说评论获取失败:", err.Error())
		return nil, err
	}
	var i int
	for range t {
		if i >= int(maxCnt) {
			break
		}
		commentsTemp, err := s.shuoShuoCommentsRaw(20, i, tid)
		if err != nil {
			log.Println("说说评论获取失败:", err.Error())
			return nil, err
		}
		if len(commentsTemp) == 0 {
			break
		}
		if len(comments) < int(num) {
			comments = append(comments, commentsTemp[0:min(len(commentsTemp), int(num)-len(comments))]...)
			i = i + 20
			time.Sleep(time.Millisecond * time.Duration(ms))
		}

	}
	return comments, nil
}

// GetLatestShuoShuo 获取用户QQ号为uin的最新说说（有空间访问权限即可）
func (s *ShuoShuoAPI) Latest(uin int64) (*ShuoShuoResp, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	ss, err := s.shuoShuoListRaw(uin, 1, 0, 0)
	if err != nil {
		er := errors.New("最新说说获取错误:" + err.Error())
		log.Println("最新说说获取失败:", er.Error())
		return nil, er
	}
	return ss[0], nil
}

// GetShuoShuoCommentsRaw 从第pos条评论开始获取num条评论，num最大为20
func (s *ShuoShuoAPI) shuoShuoCommentsRaw(num int, pos int, tid string) ([]*Comment, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	url := fmt.Sprintf(getCommentsURL, strconv.FormatInt(s.Qq, 10), pos, num, tid, s.Gtk2)
	var data string
	err := requests.
		URL(url).
		Header("cookie", s.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"cookie": s.Cookie,
	// })))
	if err != nil {
		// er := errors.New("说说评论列表请求错误:" + err.Error())
		log.Println("说说评论列表获取失败:", err)
		// return nil, err
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		// er := errors.New("说说评论正则解析错误:" + err.Error())
		log.Println("说说评论列表获取失败:", err)
		// return nil, er
	}
	jsonRaw := r[1]

	// 取出评论数据
	commentJsonList := gjson.Get(jsonRaw, "commentlist").Array()
	comments := make([]*Comment, 0)
	for _, com := range commentJsonList {
		comment := &Comment{
			ShuoShuoID: tid,
			OwnerName:  com.Get("owner.name").String(),
			OwnerUin:   com.Get("owner.uin").Int(),
			Content:    com.Get("content").String(),
			PicContent: make([]string, 0),
			CreateTime: time.Unix(com.Get("create_time").Int(), 0),
		}
		// 添加图片评论的图片到结构体
		for _, pic := range com.Get("rich_info").Array() {
			comment.PicContent = append(comment.PicContent, pic.Get("burl").String())
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

// UploadImage 上传图片
func (s *ShuoShuoAPI) uploadImage(base64img string) (*UploadImageResp, error) {
	uir := UploadImageRequest{
		Filename:      "filename",
		Uin:           s.Qq,
		Skey:          s.Skey,
		Zzpaneluin:    s.Qq,
		PUin:          s.Qq,
		PSkey:         s.Pskey,
		Uploadtype:    "1",
		Albumtype:     "7",
		Exttype:       "0",
		Refer:         "shuoshuo",
		OutputType:    "json",
		Charset:       "utf-8",
		OutputCharset: "utf-8",
		UploadHd:      "1",
		HdWidth:       "2048",
		HdHeight:      "10000",
		HdQuality:     "96",
		BackUrls:      "http://upbak.photo.qzone.qq.com/cgi-bin/upload/cgi_upload_image,http://119.147.64.75/cgi-bin/upload/cgi_upload_image",
		URL:           fmt.Sprintf(uploadImageURL, s.Gtk2),
		Base64:        "1",
		Picfile:       base64img,
		Qzreferrer:    userQzoneURL + "/" + strconv.FormatInt(s.Qq, 10),
	}

	url := fmt.Sprintf(uploadImageURL, s.Gtk2)
	// payload := strings.NewReader(structToStr(uir))
	var data string
	V := strings.NewReader(common.StructToStr(uir))
	err := requests.URL(url).
		Header("referer", userQzoneURL).
		UserAgent(common.UA).ContentType(common.ContentType).
		Header("origin", userQzoneURL).
		Header("cookie", s.Cookie).
		BodyReader(V).
		ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url), WithBody(payload),
	// 	WithHeader(map[string]string{
	// 		"referer": userQzoneURL,
	// 		"origin":  userQzoneURL,
	// 		"cookie":  s.Cookie,
	// 	})))
	if err != nil {
		log.Println("上传图片失败 ", err)
		log.Println(data)
	}
	r := cRe.FindStringSubmatch(data)
	if len(r) < 2 {
		return nil, errors.New("图片上传响应解析错误:" + data)
	}
	jsonStr := r[1]
	uploadImageResp := &UploadImageResp{
		Pre:        gjson.Get(jsonStr, "data.pre").String(),
		URL:        gjson.Get(jsonStr, "data.url").String(),
		Width:      gjson.Get(jsonStr, "data.width").Int(),
		Height:     gjson.Get(jsonStr, "data.height").Int(),
		OriginURL:  gjson.Get(jsonStr, "data.origin_url").String(),
		Contentlen: gjson.Get(jsonStr, "data.contentlen").Int(),
		Ret:        gjson.Get(jsonStr, "ret").Int(),
		Albumid:    gjson.Get(jsonStr, "data.albumid").String(),
		Lloc:       gjson.Get(jsonStr, "data.lloc").String(),
		Sloc:       gjson.Get(jsonStr, "data.sloc").String(),
		Type:       gjson.Get(jsonStr, "data.type").Int(),
	}
	return uploadImageResp, nil
}

// getPicBoAndRichval 获取已上传图片重要信息
func (s *ShuoShuoAPI) getPicBoAndRichval(data *UploadImageResp) (picBo, richval string, err error) {
	var flag bool
	if data.Ret != 0 {
		err = errors.New("已上传图片信息错误:fuck")
		return
	}
	_, picBo, flag = strings.Cut(data.URL, "&bo=")
	if !flag {
		err = errors.New("已上传图片URL错误:" + data.URL)
		return
	}
	richval = fmt.Sprintf(",%s,%s,%s,%d,%d,%d,,%d,%d", data.Albumid, data.Lloc, data.Sloc, data.Type, data.Height, data.Width, data.Height, data.Width)
	return
}

// ShuoShuoListRaw 获取用户qq号为uin且最多num个说说列表，每个说说获取上限replynum个评论数量
func (s *ShuoShuoAPI) shuoShuoListRaw(uin int64, num int, pos int, replynum int) ([]*ShuoShuoResp, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	mlr := MsgListRequest{
		Uin:                uin,
		Ftype:              "0",
		Sort:               "0",
		Pos:                strconv.Itoa(pos),
		Num:                strconv.Itoa(num),
		Replynum:           strconv.Itoa(replynum),
		GTk:                s.Gtk2,
		Callback:           "_preloadCallback",
		CodeVersion:        "1",
		Format:             "json",
		NeedPrivateComment: "1",
	}
	// url := msglistURL + structToStr(mlr)
	var data string
	url := msglistURL + common.StructToStr(mlr)
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		UserAgent(common.UA).ContentType(common.ContentType).
		Header("origin", userQzoneURL).
		Header("cookie", s.Cookie).ToString(&data).Fetch(context.Background())
	// data, err := DialRequest(NewRequest(WithUrl(url), WithHeader(map[string]string{
	// 	"referer": userQzoneURL,
	// 	"origin":  userQzoneURL,
	// 	"cookie":  s.Cookie,
	// })))
	if err != nil {
		// er := errors.New("说说列表请求错误:" + err.Error())
		log.Println("说说列表获取失败:", err)
		// return nil, er
	}
	jsonStr := data
	// 判断是否有访问权限
	forbid := gjson.Get(jsonStr, "message").String()
	if forbid != "" {
		er := errors.New("说说列表解析错误:" + forbid)
		log.Println("说说列表获取失败:", er.Error())
		return nil, er
	}

	var resLen int64
	if !gjson.Get(jsonStr, "msglist.#").Exists() {
		er := errors.New("说说列表解析错误:" + jsonStr)
		log.Println("说说列表获取失败:", er.Error())
		return nil, er
	}
	resLen = gjson.Get(jsonStr, "msglist.#").Int()
	results := make([]*ShuoShuoResp, min(resLen, int64(num)))
	index := 0

	lists := gjson.Get(jsonStr, "msglist").Array()
	for _, shuoshuo := range lists {
		ss := &ShuoShuoResp{
			Uin:         shuoshuo.Get("uin").Int(),
			Name:        shuoshuo.Get("name").String(),
			Tid:         shuoshuo.Get("tid").String(),
			Content:     shuoshuo.Get("content").String(),
			CreateTime:  shuoshuo.Get("createTime").String(),
			CreatedTime: shuoshuo.Get("created_time").Int(),
			PicTotal:    shuoshuo.Get("pictotal").Int(),
			Cmtnum:      shuoshuo.Get("cmtnum").Int(),
			Secret:      shuoshuo.Get("secret").Int(),
		}

		pics := shuoshuo.Get("pic").Array()
		for _, pic := range pics {
			ss.Pic = append(ss.Pic, PicResp{
				PicId:      pic.Get("pic_id").String(),
				Url1:       pic.Get("url1").String(),
				Url2:       pic.Get("url2").String(),
				Url3:       pic.Get("url3").String(),
				Smallurl:   pic.Get("smallurl").String(),
				Curlikekey: pic.Get("curlikekey").String(),
			})
		}

		results[index] = ss
		index++
	}
	return results, nil
}

// DoLike 说说空间点赞 TODO:疑似无效
func (s *ShuoShuoAPI) DoLike(tid string) (*LikeResp, error) {
	if s.Err != nil {
		return nil, s.Err
	}
	lr := LikeRequest{
		Qzreferrer: userQzoneURL + strconv.FormatInt(s.Qq, 10),
		Opuin:      s.Qq,
		Unikey:     userQzoneURL + strconv.FormatInt(s.Qq, 10) + "/mood/" + tid,
		From:       "1",
		Fid:        tid,
		Typeid:     "0",
		Appid:      "311",
	}
	lr.Curkey = lr.Unikey
	url := fmt.Sprintf(likeURL, s.Gtk2)
	var data string
	V := strings.NewReader(common.StructToStr(lr))
	err := requests.
		URL(url).
		Header("referer", userQzoneURL).
		UserAgent(common.UA).ContentType(common.ContentType).
		Header("origin", userQzoneURL).
		Header("cookie", s.Cookie).
		BodyReader(V).ToString(&data).Fetch(context.Background())
	// payload := strings.NewReader(structToStr(lr))
	// data, err := DialRequest(NewRequest(WithMethod("POST"), WithUrl(url),
	// 	WithBody(payload), WithHeader(map[string]string{
	// 		"referer": userQzoneURL,
	// 		"origin":  userQzoneURL,
	// 		"cookie":  s.Cookie,
	// 	})))
	if err != nil {
		// er := errors.New("点赞请求错误:" + err.Error())
		log.Println("空间点赞失败:", err)
		// return nil, er
	}
	r := cReLike.FindStringSubmatch(data)
	if len(r) < 2 {
		er := errors.New("点赞响应解析错误:" + data)
		log.Println("空间点赞失败:", er.Error())
		return nil, er
	}
	likeResp := &LikeResp{
		Ret: gjson.Get(r[1], "ret").Int(),
		Msg: gjson.Get(r[1], "msg").String(),
	}
	if likeResp.Msg != "succ" {
		er := errors.New("点赞未生效" + likeResp.Msg)
		log.Println("空间点赞失败:", er.Error())
		return nil, er
	}
	return likeResp, nil
}
