package test

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dop251/goja"
	"github.com/intchensc/qzone"
	"github.com/intchensc/qzone/auth"
	"github.com/tidwall/gjson"
)

func TestFriendAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试好友列表
	friends, err := qm.API.Friend().List()
	if err != nil {
		t.Fatal("获取好友列表失败:", err)
	}
	t.Logf("获取到%d个好友", len(friends))

	// 测试好友详情
	if len(friends) > 0 {
		_, err := qm.API.Friend().Detail(friends[0].Uin)
		if err != nil {
			t.Fatal("获取好友详情失败:", err)
		}
	}
}

func TestGroupAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试群列表
	groups, err := qm.API.Group().List()
	if err != nil {
		t.Fatal("获取群列表失败:", err)
	}
	t.Logf("获取到%d个群", len(groups))

	// 测试群成员
	if len(groups) > 0 {
		_, err := qm.API.Group().MemberList(groups[0].GroupCode)
		if err != nil {
			t.Fatal("获取群成员失败:", err)
		}
	}
}

// TODO: 重写历史记录
func TestMessageAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}
	fmt.Println(qm.API.Cookie)
	// 测试历史消息
	history, err := qm.API.Message().ListN(19, 10)
	if err != nil {
		t.Fatal("获取历史消息失败:", err)
	}
	for _, k := range history {
		t.Logf(k.TID, k.QQ, k.Type, k.CreateTime)
	}
	t.Logf("获取指定个数，获取到%d条历史消息", len(history))

	historyT, err := qm.API.Message().ListT(time.Unix(1758381345, 0), 10)
	for _, k := range historyT {
		t.Logf(k.TID, k.QQ, k.Type, k.CreateTime)
	}
	if err != nil {
		t.Fatal("获取历史消息失败:", err)
	}
	t.Logf("获取指定个数，获取到%d条历史消息", len(historyT))
}

func TestShuoShuoAPI(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 测试说说数量
	count, err := qm.API.ShuoShuo().Count(qm.API.Qq)
	if err != nil {
		t.Fatal("获取说说数量失败:", err)
	}
	t.Logf("说说总数:%d", count)

	// 测试说说列表
	if count > 0 {
		ss, err := qm.API.ShuoShuo().List(qm.API.Qq, 3, 1000)
		if err != nil {
			t.Fatal("获取说说列表失败:", err)
		}
		t.Logf("获取到说说:%v", ss[0].Content)
	}

	// 测试最新说说
	if count > 0 {
		_, err := qm.API.ShuoShuo().Latest(qm.API.Qq)
		if err != nil {
			t.Fatal("获取最新说说失败:", err)
		}
	}

	// 测试说说评论
	if count > 0 {
		list, err := qm.API.ShuoShuo().List(qm.API.Qq, 1, 0)
		if err != nil {
			t.Fatal("获取说说列表失败:", err)
		}
		if len(list) > 0 {
			_, err := qm.API.ShuoShuo().CommentList(list[0].Tid, 10, 0)
			if err != nil {
				t.Fatal("获取说说评论失败:", err)
			}
		}
	}
}

// 测试发布纯文本说说
func TestPublishTextShuoShuo(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	resp, err := qm.API.ShuoShuo().Publish("测试纯文本说说", nil)
	if err != nil {
		t.Fatal("发布纯文本说说失败:", err)
	}
	t.Logf("发布成功, tid:%v", resp)
}

// 测试发布带图片说说
func TestPublishImageShuoShuo(t *testing.T) {
	qm := qzone.New(&auth.QrAuth{})
	if err := qm.Login(); err != nil {
		t.Fatal("登录失败:", err)
	}

	// 从test目录读取图片数据
	imagePath := "./1.png"
	imageData, err := os.ReadFile(imagePath)
	if err != nil {
		t.Fatal("读取图片文件失败:", err)
	}
	images := []string{base64.StdEncoding.EncodeToString(imageData)}
	resp, err := qm.API.ShuoShuo().Publish("测试带图片说说", images)
	if err != nil {
		t.Fatal("发布带图片说说失败:", err)
	}
	t.Logf("发布成功, tid:%v", resp)
}

func TestName(t *testing.T) {
	jsonRaw := "{main:{hasMoreFeeds:true,total_number:3,offset:0,count:2,begin_time:'',end_time:'',getappnotification:'',has_get_key:'',host_more:''},data:[{ver:'1',appid:'311',typeid:'4',key:'4844244dee06cd68a6c10000',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'0',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267117',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'                                                                                                                                                             \\x3Cli class=\\x22f-single f-s-s\\x22 id=\\x22fct_1294222408_311_4_1758267117_1_1\\x22>\\x3Cdiv class=\\x22f-single-head f-aside\\x22>\\x3Cdiv class=\\x22f-adorn-top\\x22>\\x3C\\/div>    \\x3Cdiv class=\\x22user-pto\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22 target=\\x22_blank\\x22 class=\\x22user-avatar q_namecard f-s-a\\x22 link=\\x22nameCard_1294222408\\x22 data-clicklog=\\x22avatar\\x22>\\x3Cimg src=\\x22https:\\/\\/qlogo1.store.qq.com\\/qzone\\/1294222408\\/1294222408\\/50?1687609779\\x22>\\x3C\\/a>\\x3C\\/div>\\x3Cdiv class=\\x22user-op\\x22>\\x3Ca href=\\x22javascript:;\\x22 class=\\x22arrow-down\\x22 data-cmd=\\x22qz_opr_more\\x22 data-moreoperate=\\x220\\x22>\\x3Ci class=\\x22fui-icon icon-arrow-down\\x22>\\x3C\\/i>\\x3C\\/a>        \\x3C\\/div>\\x3Cdiv class=\\x22user-info\\x22>\\x3Cdiv class=\\x22f-nick\\x22>\\x3Ca target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22  data-clicklog=\\x22nick\\x22 class=\\x22f-name q_namecard \\x22 link=\\x22nameCard_1294222408\\x22>R\\x3C\\/a>  \\x3Cspan  class=\\x22 ui-mr10 state\\x22 >提到我\\x3C\\/span>\\x3C\\/div>\\x3Cdiv class=\\x22info-detail\\x22>\\x3Cspan  class=\\x22 ui-mr8 state\\x22 >昨天 15:31\\x3C\\/span>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-content f-wrap\\x22>\\x3Cdiv class=\\x22f-item f-s-i f-item-passive\\x22 id=\\x22feed_1294222408_311_4_1758267117_1_1\\x22 data-feedsflag=\\x22\\x22 data-iswupfeed=\\x221\\x22 data-key=\\x224844244dee06cd68a6c10000\\x22 data-specialtype=\\x22\\x22 data-extend-info=\\x220_0_0_0_0_0_0|18c5950002028003|0008000056d05c05\\x22 data-functype=\\x22\\x22 data-hasfollowed=\\x221\\x22>\\x3Cdiv class=\\x22qz_summary wupfeed\\x22 id=\\x22hex_1294222408_311_4_1758267117_1_1\\x22>\\x3Ci class=\\x22none\\x22 name=\\x22feed_data\\x22  data-bitmap=\\x2218c5950002028003\\x22 data-yybitmap=\\x220008000056d05c05\\x22 data-vipstarbitmap=\\x221010088008000080\\x22 data-fkey=\\x224844244dee06cd68a6c10000\\x22 data-tid=\\x224844244dee06cd68a6c10000\\x22 data-uin=\\x221294222408\\x22 data-origfkey=\\x22\\x22 data-origtid=\\x224844244d5104cd6851410300\\x22 data-origuin=\\x221294222408\\x22 data-subid=\\x22\\x22 data-totweet=\\x22\\x22 data-issignin=\\x22\\x22 data-source=\\x22\\x22 data-retweetcount=\\x220\\x22 data-stat=\\x22XkyBMA0tUiPDkbJUMpEkCseurzFFTpPzZenkOpw\\/2aeAf6Qv6eSCgcBD!7r0L3XOljB8iPtRfA2SriLVetGC9QE5lu4P\\/fv\\/XsRRPEjk\\/ULF6EOZeLEnDdLfdmlNMuUnjYISrhsp9AU_\\x22   data-topicid=\\x221294222408_4844244dee06cd68a6c10000__1\\x22    data-feedstype=\\x22103\\x22  data-abstime=\\x221758267117\\x22   data-iswupfeed=\\x221\\x22  data-platformid=\\x2252\\x22  data-accessright=\\x221\\x22>\\x3C\\/i>\\x3Cdiv class=\\x22f-ct \\x22 data-detailurl=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244dee06cd68a6c10000.\\x22>\\x3Cdiv class=\\x22fui-left-right f-ct-txtimg f-ct-imgtxt f-ct-b-np\\x22>\\x3Cdiv class=\\x22img-box\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244dee06cd68a6c10000.\\x22 target=\\x22_blank\\x22>                        \\x3Ci class=\\x22ui-icon icon-default-shuoshuo\\x22>\\x3C\\/i>                                \\x3C\\/a>                            \\x3C\\/div>\\x3Cdiv class=\\x22inner\\x22>\\x3Cdiv>\\x3Cdiv class=\\x22txt-box one-line\\x22>\\x3Cdiv class=\\x22inner\\x22>\\x3Cp class=\\x22txt-box-title ellipsis-one\\x22>\\x3Cspan class=\\x22ellipsis-front\\x22>\\x3Ca class=\\x22nickname name c_tx ui-mr5 q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp;\\x3Cspan  class=\\x22 ui-mr5 state\\x22 >：\\x3C\\/span>这是转发艾特人带评论\\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_941094692\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/941094692\\x22>@圈\\x3C\\/a>&nbsp;\\x3C\\/span>\\x3Cspan class=\\x22state\\x22> || \\x3C\\/span> \\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp;\\x3Cspan  class=\\x22 state\\x22 >：\\x3C\\/span>这是一条测试说说\\x3C\\/p>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-foot\\x22>        \\x3Cdiv class=\\x22mod-comments\\x22 style=\\x22padding:0\\x22>\\x3Cdiv class=\\x22mod-commnets-poster feedClickCmd comment_default_inputentry\\x22 data-cmd=\\x22qz_reply\\x22 data-version=\\x226\\x22 data-action=\\x22http:\\/\\/taotao.qq.com\\/cgi-bin\\/emotion_cgi_re_feeds\\x22 data-param=\\x22t1_source=&amp;t1_uin=1294222408&amp;t1_tid=4844244dee06cd68a6c10000&amp;signin=&amp;sceneid=103\\x22 data-charset=\\x22utf-8\\x22 data-maxLength=\\x22\\x22 data-tuin=\\x221294222408\\x22 data-config=\\x221|1|1|1,b52,with_fwd,同时转发;0|1,taotaoact.qzone.qq.com,@InputReply|1,taotaoact.qzone.qq.com,@ClickReply|1,taotaoact.qzone.qq.com,commentPresentClick\\x22  data-tid=\\x22\\x22 >\\x3Cdiv class=\\x22comments-poster-bd comments-poster-default\\x22>\\x3Cdiv class=\\x22comments-box\\x22 data-clicklog=\\x22comment\\x22>\\x3Cdiv class=\\x22textinput textinput-default bor2\\x22 contenteditable=\\x22true\\x22 alt=\\x22replybtn\\x22 placeholder=\\x22评论\\x22>\\x3Ca class=\\x22c_tx3\\x22 href=\\x22javascript:void(0);\\x22 alt=\\x22replybtn\\x22>评论\\x3C\\/a>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>    \\x3C\\/div>\\x3C\\/li>',opuin:'1294222408'},{ver:'1',appid:'311',typeid:'4',key:'4_4844244d5104cd6851410300',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'1',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267117',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'                                                                                                                                                                 \\x3Cli class=\\x22f-single f-s-s\\x22 id=\\x22fct_1294222408_311_4_1758267117_1_1\\x22>\\x3Cdiv class=\\x22f-single-head f-aside\\x22>\\x3Cdiv class=\\x22f-adorn-top\\x22>\\x3C\\/div>    \\x3Cdiv class=\\x22user-pto\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22 target=\\x22_blank\\x22 class=\\x22user-avatar q_namecard f-s-a\\x22 link=\\x22nameCard_1294222408\\x22 data-clicklog=\\x22avatar\\x22>\\x3Cimg src=\\x22https:\\/\\/qlogo1.store.qq.com\\/qzone\\/1294222408\\/1294222408\\/50?1687609779\\x22>\\x3C\\/a>\\x3C\\/div>\\x3Cdiv class=\\x22user-op\\x22>\\x3Ca href=\\x22javascript:;\\x22 class=\\x22arrow-down\\x22 data-cmd=\\x22qz_opr_more\\x22 data-moreoperate=\\x220\\x22>\\x3Ci class=\\x22fui-icon icon-arrow-down\\x22>\\x3C\\/i>\\x3C\\/a>        \\x3C\\/div>\\x3Cdiv class=\\x22user-info\\x22>\\x3Cdiv class=\\x22f-nick\\x22>\\x3Ca target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22  data-clicklog=\\x22nick\\x22 class=\\x22f-name q_namecard \\x22 link=\\x22nameCard_1294222408\\x22>R\\x3C\\/a>  \\x3Cspan  class=\\x22 ui-mr10 state\\x22 >评论提到我\\x3C\\/span>\\x3C\\/div>\\x3Cdiv class=\\x22info-detail\\x22>\\x3Cspan  class=\\x22 ui-mr8 state\\x22 >昨天 15:31\\x3C\\/span>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-content f-wrap\\x22>\\x3Cdiv class=\\x22f-item f-s-i f-item-passive\\x22 id=\\x22feed_1294222408_311_4_1758267117_1_1\\x22 data-feedsflag=\\x22\\x22 data-iswupfeed=\\x221\\x22 data-key=\\x224_4844244d5104cd6851410300\\x22 data-specialtype=\\x22\\x22 data-extend-info=\\x220_0_0_0_0_0_0|18c5950002028003|0008000056d05c05\\x22 data-functype=\\x22\\x22 data-hasfollowed=\\x221\\x22>\\x3Cdiv class=\\x22qz_summary wupfeed\\x22 id=\\x22hex_1294222408_311_4_1758267117_1_1\\x22>\\x3Ci class=\\x22none\\x22 name=\\x22feed_data\\x22  data-bitmap=\\x2218c5950002028003\\x22 data-yybitmap=\\x220008000056d05c05\\x22 data-vipstarbitmap=\\x221010088008000080\\x22 data-fkey=\\x224_4844244d5104cd6851410300\\x22 data-tid=\\x224844244d5104cd6851410300\\x22 data-uin=\\x221294222408\\x22 data-origfkey=\\x22\\x22 data-origtid=\\x224844244d5104cd6851410300\\x22 data-origuin=\\x221294222408\\x22 data-subid=\\x22\\x22 data-totweet=\\x22\\x22 data-issignin=\\x22\\x22 data-source=\\x22\\x22 data-retweetcount=\\x220\\x22 data-stat=\\x22XkyBMA0tUiPDkbJUMpEkCseurzFFTpPzt4e69jcjVb8oEnxWDOVRK7u\\/xEsUYJ4oljB8iPtRfA2SriLVetGC9QE5lu4P\\/fv\\/XsRRPEjk\\/ULF6EOZeLEnDbrAzp37QLSuYbfWt3I5eIQ_\\x22   data-topicid=\\x221294222408_4844244d5104cd6851410300__1\\x22    data-feedstype=\\x22107\\x22  data-abstime=\\x221758267117\\x22   data-iswupfeed=\\x221\\x22  data-platformid=\\x2252\\x22  data-accessright=\\x221\\x22>\\x3C\\/i>\\x3Cdiv class=\\x22f-ct \\x22 data-detailurl=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244d5104cd6851410300.\\x22>\\x3Cdiv class=\\x22fui-left-right f-ct-txtimg f-ct-imgtxt f-ct-b-np\\x22>\\x3Cdiv class=\\x22img-box\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244d5104cd6851410300.\\x22 target=\\x22_blank\\x22>                        \\x3Ci class=\\x22ui-icon icon-default-shuoshuo\\x22>\\x3C\\/i>                                \\x3C\\/a>                            \\x3C\\/div>\\x3Cdiv class=\\x22txt-box one-line\\x22>\\x3Cdiv class=\\x22inner\\x22>\\x3Cp class=\\x22txt-box-title ellipsis-one\\x22>\\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp;\\x3Cspan  class=\\x22 state\\x22 >：\\x3C\\/span>这是一条测试说说\\x3C\\/p>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-foot\\x22>        \\x3Cdiv class=\\x22f-ang-t\\x22>\\x3C\\/div>\\x3Cdiv class=\\x22mod-comments\\x22 style=\\x22padding:0\\x22>\\x3Cdiv class=\\x22comments-list \\x22>\\x3Cul>\\x3Cli class=\\x22comments-item bor3\\x22 data-type=\\x22commentroot\\x22 data-tid=\\x224\\x22 data-uin=\\x221294222408\\x22 data-nick=\\x22R\\x22 data-who=\\x221\\x22>\\x3Cdiv class=\\x22comments-item-bd\\x22>\\x3Cdiv class=\\x22single-reply\\x22>\\x3Cdiv class=\\x22ui-avatar\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22 target=\\x22_blank\\x22>\\x3Cimg class=\\x22q_namecard\\x22 link=\\x22nameCard_1294222408 des_1294222408\\x22 alt=\\x22R\\x22 src=\\x22http:\\/\\/qlogo1.store.qq.com\\/qzone\\/1294222408\\/1294222408\\/30?1687609779\\x22 \\/>\\x3C\\/a>\\x3C\\/div>\\x3Cdiv class=\\x22comments-content font-b\\x22>\\x3Ca class=\\x22nickname name c_tx q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp; : 这是转发艾特人带评论\\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_941094692\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/941094692\\x22>@圈\\x3C\\/a>&nbsp;\\x3Cdiv class=\\x22comments-op\\x22>\\x3Cspan  class=\\x22 ui-mr10 state\\x22 >昨天 15:31\\x3C\\/span>\\x3Ca class=\\x22act-reply none\\x22 href=\\x22javascript:;\\x22 data-cmd=\\x22qz_reply\\x22 data-version=\\x226.4\\x22 data-action=\\x22http:\\/\\/taotao.qq.com\\/cgi-bin\\/emotion_cgi_re_feeds\\x22 data-param=\\x22t1_source=&amp;t1_uin=1294222408&amp;t1_tid=4844244d5104cd6851410300&amp;t2_uin=1294222408&amp;t2_tid=4&amp;subdotype=55702&amp;signin=0&amp;sceneid=107\\x22 data-charset=\\x22utf-8\\x22 data-tuin=\\x22\\x22 data-config=\\x221|1|1|0|1,taotaoact.qzone.qq.com,@InputReply|1,taotaoact.qzone.qq.com,@ClickReply|1,taotaoact.qzone.qq.com,commentPresentClick\\x22>\\x3Cb class=\\x22hide-clip\\x22>回复\\x3C\\/b>\\x3C\\/a>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/li>\\x3C\\/ul>\\x3C\\/div>\\x3Cdiv class=\\x22mod-commnets-poster feedClickCmd comment_default_inputentry\\x22 data-cmd=\\x22qz_reply\\x22 data-version=\\x226\\x22 data-action=\\x22http:\\/\\/taotao.qq.com\\/cgi-bin\\/emotion_cgi_re_feeds\\x22 data-param=\\x22t1_source=&amp;t1_uin=1294222408&amp;t1_tid=4844244d5104cd6851410300&amp;t2_uin=1294222408&amp;t2_tid=4&amp;subdotype=55702&amp;signin=0&amp;sceneid=107\\x22 data-charset=\\x22utf-8\\x22 data-maxLength=\\x22\\x22 data-tuin=\\x221294222408\\x22 data-config=\\x22\\x22   data-typeid=\\x221\\x22   data-tid=\\x224\\x22  data-uin=\\x221294222408\\x22  data-nick=\\x22R\\x22  data-who=\\x220\\x22 >\\x3Cdiv class=\\x22comments-poster-bd comments-poster-default\\x22>\\x3Cdiv class=\\x22comments-box\\x22 data-clicklog=\\x22comment\\x22>\\x3Cdiv class=\\x22textinput textinput-default bor2\\x22 contenteditable=\\x22true\\x22 alt=\\x22replybtn\\x22 placeholder=\\x22评论\\x22>\\x3Ca class=\\x22c_tx3\\x22 href=\\x22javascript:void(0);\\x22 alt=\\x22replybtn\\x22>评论\\x3C\\/a>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>    \\x3C\\/div>\\x3C\\/li>',opuin:'1294222408'},{ver:'1',appid:'311',typeid:'4',key:'4844244dc206cd6866090500',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'2',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267074',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'                                                                                                                                                             \\x3Cli class=\\x22f-single f-s-s\\x22 id=\\x22fct_1294222408_311_4_1758267074_1_1\\x22>\\x3Cdiv class=\\x22f-single-head f-aside\\x22>\\x3Cdiv class=\\x22f-adorn-top\\x22>\\x3C\\/div>    \\x3Cdiv class=\\x22user-pto\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22 target=\\x22_blank\\x22 class=\\x22user-avatar q_namecard f-s-a\\x22 link=\\x22nameCard_1294222408\\x22 data-clicklog=\\x22avatar\\x22>\\x3Cimg src=\\x22https:\\/\\/qlogo1.store.qq.com\\/qzone\\/1294222408\\/1294222408\\/50?1687609779\\x22>\\x3C\\/a>\\x3C\\/div>\\x3Cdiv class=\\x22user-op\\x22>\\x3Ca href=\\x22javascript:;\\x22 class=\\x22arrow-down\\x22 data-cmd=\\x22qz_opr_more\\x22 data-moreoperate=\\x220\\x22>\\x3Ci class=\\x22fui-icon icon-arrow-down\\x22>\\x3C\\/i>\\x3C\\/a>        \\x3C\\/div>\\x3Cdiv class=\\x22user-info\\x22>\\x3Cdiv class=\\x22f-nick\\x22>\\x3Ca target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22  data-clicklog=\\x22nick\\x22 class=\\x22f-name q_namecard \\x22 link=\\x22nameCard_1294222408\\x22>R\\x3C\\/a>  \\x3Cspan  class=\\x22 ui-mr10 state\\x22 >提到我\\x3C\\/span>\\x3C\\/div>\\x3Cdiv class=\\x22info-detail\\x22>\\x3Cspan  class=\\x22 ui-mr8 state\\x22 >昨天 15:31\\x3C\\/span>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-content f-wrap\\x22>\\x3Cdiv class=\\x22f-item f-s-i f-item-passive\\x22 id=\\x22feed_1294222408_311_4_1758267074_1_1\\x22 data-feedsflag=\\x22\\x22 data-iswupfeed=\\x221\\x22 data-key=\\x224844244dc206cd6866090500\\x22 data-specialtype=\\x22\\x22 data-extend-info=\\x220_0_0_0_0_0_0|18c5950002028003|0008000056d05c05\\x22 data-functype=\\x22\\x22 data-hasfollowed=\\x221\\x22>\\x3Cdiv class=\\x22qz_summary wupfeed\\x22 id=\\x22hex_1294222408_311_4_1758267074_1_1\\x22>\\x3Ci class=\\x22none\\x22 name=\\x22feed_data\\x22  data-bitmap=\\x2218c5950002028003\\x22 data-yybitmap=\\x220008000056d05c05\\x22 data-vipstarbitmap=\\x221010088008000080\\x22 data-fkey=\\x224844244dc206cd6866090500\\x22 data-tid=\\x224844244dc206cd6866090500\\x22 data-uin=\\x221294222408\\x22 data-origfkey=\\x22\\x22 data-origtid=\\x224844244d5104cd6851410300\\x22 data-origuin=\\x221294222408\\x22 data-subid=\\x22\\x22 data-totweet=\\x22\\x22 data-issignin=\\x22\\x22 data-source=\\x22\\x22 data-retweetcount=\\x220\\x22 data-stat=\\x22XkyBMA0tUiPDkbJUMpEkCseurzFFTpPzav549t03pW44LARvACXOsySg2ed4p1UcljB8iPtRfA2SriLVetGC9QE5lu4P\\/fv\\/XsRRPEjk\\/ULF6EOZeLEnDRFW2\\/9asEeNPRd\\/yxrSeXc_\\x22   data-topicid=\\x221294222408_4844244dc206cd6866090500__1\\x22    data-feedstype=\\x22103\\x22  data-abstime=\\x221758267074\\x22   data-iswupfeed=\\x221\\x22  data-platformid=\\x2252\\x22  data-accessright=\\x221\\x22>\\x3C\\/i>\\x3Cdiv class=\\x22f-ct \\x22 data-detailurl=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244dc206cd6866090500.\\x22>\\x3Cdiv class=\\x22fui-left-right f-ct-txtimg f-ct-imgtxt f-ct-b-np\\x22>\\x3Cdiv class=\\x22img-box\\x22>\\x3Ca href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\/mood\\/4844244dc206cd6866090500.\\x22 target=\\x22_blank\\x22>                        \\x3Ci class=\\x22ui-icon icon-default-shuoshuo\\x22>\\x3C\\/i>                                \\x3C\\/a>                            \\x3C\\/div>\\x3Cdiv class=\\x22inner\\x22>\\x3Cdiv>\\x3Cdiv class=\\x22txt-box one-line\\x22>\\x3Cdiv class=\\x22inner\\x22>\\x3Cp class=\\x22txt-box-title ellipsis-one\\x22>\\x3Cspan class=\\x22ellipsis-front\\x22>\\x3Ca class=\\x22nickname name c_tx ui-mr5 q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp;\\x3Cspan  class=\\x22 ui-mr5 state\\x22 >：\\x3C\\/span>这是转发艾特人不带评论\\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_941094692\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/941094692\\x22>@圈\\x3C\\/a>&nbsp;\\x3C\\/span>\\x3Cspan class=\\x22state\\x22> || \\x3C\\/span> \\x3Ca class=\\x22nickname name c_tx  q_namecard\\x22  link=\\x22nameCard_1294222408\\x22 target=\\x22_blank\\x22 href=\\x22http:\\/\\/user.qzone.qq.com\\/1294222408\\x22>R\\x3C\\/a>&nbsp;\\x3Cspan  class=\\x22 state\\x22 >：\\x3C\\/span>这是一条测试说说\\x3C\\/p>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3Cdiv class=\\x22f-single-foot\\x22>        \\x3Cdiv class=\\x22mod-comments\\x22 style=\\x22padding:0\\x22>\\x3Cdiv class=\\x22mod-commnets-poster feedClickCmd comment_default_inputentry\\x22 data-cmd=\\x22qz_reply\\x22 data-version=\\x226\\x22 data-action=\\x22http:\\/\\/taotao.qq.com\\/cgi-bin\\/emotion_cgi_re_feeds\\x22 data-param=\\x22t1_source=&amp;t1_uin=1294222408&amp;t1_tid=4844244dc206cd6866090500&amp;signin=&amp;sceneid=103\\x22 data-charset=\\x22utf-8\\x22 data-maxLength=\\x22\\x22 data-tuin=\\x221294222408\\x22 data-config=\\x221|1|1|1,b52,with_fwd,同时转发;0|1,taotaoact.qzone.qq.com,@InputReply|1,taotaoact.qzone.qq.com,@ClickReply|1,taotaoact.qzone.qq.com,commentPresentClick\\x22  data-tid=\\x22\\x22 >\\x3Cdiv class=\\x22comments-poster-bd comments-poster-default\\x22>\\x3Cdiv class=\\x22comments-box\\x22 data-clicklog=\\x22comment\\x22>\\x3Cdiv class=\\x22textinput textinput-default bor2\\x22 contenteditable=\\x22true\\x22 alt=\\x22replybtn\\x22 placeholder=\\x22评论\\x22>\\x3Ca class=\\x22c_tx3\\x22 href=\\x22javascript:void(0);\\x22 alt=\\x22replybtn\\x22>评论\\x3C\\/a>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>\\x3C\\/div>    \\x3C\\/div>\\x3C\\/li>',opuin:'1294222408'},undefined]}\n"
	//jsonRawa := "{main:{hasMoreFeeds:true,total_number:3,offset:0,count:2,begin_time:'',end_time:'',getappnotification:'',has_get_key:'',host_more:''},data:[{ver:'1',appid:'311',typeid:'4',key:'4844244dee06cd68a6c10000',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'0',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267117',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'',opuin:'1294222408'},{ver:'1',appid:'311',typeid:'4',key:'4_4844244d5104cd6851410300',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'1',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267117',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'',opuin:'1294222408'},{ver:'1',appid:'311',typeid:'4',key:'4844244dc206cd6866090500',flag:'1',dataonly:'0',titleTemp:'',summaryTemp:'',feedno:'2',title:'',summary:'',appiconid:'311',clscFold:'icenter_list_extend',abstime:'1758267074',feedstime:'昨天 15:31',userHome:'http://user.qzone.qq.com/1294222408',namecardLink:'nameCard_1294222408',ouin:'',uin:'1294222408',foldFeed:'',foldFeedTitle:'',showEbtn:'',scope:'1',hideExtend:'',nickname:'R',emoji:[],remark:'',type:'',vip:'vip_0',bitmap:'18c5950002028003',             yybitmap:'0008000056d05c05',info_user_name:'',logimg:'http://qlogo1.store.qq.com/qzone/1294222408/1294222408/50?1687609779',bor:'',lastFeedBor:'',list_bor2:'',info_user_display:'',upernum:'',oprType:'0',moreflag:'',otherflag:'0_0_0_0_0_0_0',rightflag:'',sameuser:{},uper_isfriend:[],uperlist:[],html:'',opuin:'1294222408'},undefined]}"
	//j := `{"main":{"a":1,"b":2},"data":{"data":[]}}`

	vm := goja.New()
	// 在 JS 环境中，直接 eval 就能处理这种写法
	result, err := vm.RunString(fmt.Sprintf("JSON.stringify(%s)", jsonRaw))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(gjson.Get(result.String(), "data").String())
}
