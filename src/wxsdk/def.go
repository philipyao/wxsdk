package wxsdk

import (
	"net/http"
	"encoding/xml"
)

const (
	Token               = "keepmovingbuddy"
)

type MsgType string
const (
	MsgTypeEvent                MsgType = "event"
    MsgTypeText                 MsgType = "text"
	MsgTypeImage                MsgType = "image"
	MsgTypeVoice                MsgType = "voice"
	MsgTypeVideo                MsgType = "video"
	MsgTypeShortVideo           MsgType = "shortvideo"
	MsgTypeLink                 MsgType = "link"

	EventTypeSubscribe          MsgType = "event.subscribe"
	EventTypeUnsubscribe        MsgType = "event.unsubscribe"
)

const (
	MediaTypeImage              = "image"
	MediaTypeVoice              = "voice"   //语音
	MediaTypeVideo              = "video"
	MediaTypeThumb              = "thumb"   //缩略图

	//临时
	UrlTempMediaUpload          = "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%v&type=%v"
	UrlTempMediaGet             = "https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"

	//永久
	UrlMediaUpload              = "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%v&type=%v"
	UrlMediaUploadNews          = "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=%v"
	UrlMediaUploadNewsImg       = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%v"
)

const (
	MenuTypeClick               = "click"              // 点击推事件
	MenuTypeView                = "view"               // 点击推事件
	MenuTypeScancodePush        = "scancode_push"      // 扫码推事件
	MenuTypeScancodeWaitmsg     = "scancode_waitmsg"   // 扫码推事件且弹出“消息接收中”提示框
	MenuTypePicSysphoto         = "pic_sysphoto"       // 弹出系统拍照发图
	MenuTypePicPhotoOrAlbum     = "pic_photo_or_album" // 弹出拍照或者相册发图
	MenuTypePicWeixin           = "pic_weixin"         // 弹出微信相册发图器
	MenuTypeLocationSelect      = "location_select"    // 弹出地理位置选择器
	MenuTypeMediaId             = "media_id"           // 下发消息（除文本消息）
	MenuTypeViewLimited         = "view_limited"       // 跳转图文消息URL

	UrlMenuCreate               = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%v"
	UrlMenuGet                  = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%v"
	UrlMenuDelete               = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%v"
)

type MsgBase struct {
	XMLName         xml.Name `xml:"xml"`
	ToUserName      string   `xml:"ToUserName"`             //开发者微信号
	FromUserName    string   `xml:"FromUserName"`           //发送方帐号（一个OpenID）
	CreateTime      uint32   `xml:"CreateTime"`             //消息创建时间
	MsgType         MsgType  `xml:"MsgType"`                //消息类型
}

type MsgPlain struct {
	MsgId           uint64   `xml:"MsgId,omitempty"`          //消息id，用于去重
	Content         string   `xml:"Content,omitempty"`        //文本消息内容
	PicUrl          string   `xml:"PicUrl,omitempty"`         //图片链接
	Format          string   `xml:"Format,omitempty"`         //语音格式
	MediaId         string   `xml:"MediaId,omitempty"`        //图片媒体id
	ThumbMediaId    string   `xml:"ThumbMediaId,omitempty"`   //缩略图
	Location_X      float32  `xml:"Location_X,omitempty"`     //地理位置纬度
	Location_Y      float32  `xml:"Location_Y,omitempty"`     //地理位置经度
	Scale           uint32   `xml:"Scale,omitempty"`          //地图缩放大小
	Label           string   `xml:"Label,omitempty"`          //地理位置信息
	Title           string   `xml:"Title,omitempty"`          //消息标题
	Description     string   `xml:"Description,omitempty"`    //消息描述
	Url             string   `xml:"Url,omitempty"`            //消息链接
}

type MsgEvent struct {
	Event           string   `xml:"Event,omitempty"`          //事件类型
	EventKey        string   `xml:"EventKey,omitempty"`       //事件KEY值
	Ticket          string   `xml:"Ticket,omitempty"`         //二维码的ticket
	Latitude        float32  `xml:"Latitude,omitempty"`       //地理位置纬度
	Longitude       float32  `xml:"Longitude,omitempty"`      //地理位置经度
	Precision       float32  `xml:"Precision,omitempty"`      //地理位置精度
}

//微信服务器发过来的消息
type Message struct {
	MsgBase

	MsgPlain

	MsgEvent
}

//请求上下文
type RequestContext struct {
	MsgBase
	w http.ResponseWriter
}

//请求微信接口，通用返回
type WeiXinRspHeader struct {
	ErrorCode    int    `json:"errcode,omitempty"`
	ErrorMessage string `json:"errmsg,omitempty"`
}
