package wxsdk

import(
	//"time"
	"fmt"
	"io/ioutil"
	"crypto/sha1"
	"sort"
	"strings"
    "net/http"
	"encoding/xml"

    "wxsdk/wxproto"
)

type WXSDKHandle struct {}

var (
	defaultServeMux WXSDKHandle
	handles map[wxproto.MsgType]handleFunc = make(map[wxproto.MsgType]handleFunc)
)

type handleFunc func(ctx *wxproto.RequestContext, reqMsg *wxproto.Message)

func init() {
	handles[wxproto.MsgTypeText]            = defaultPlainHandle
	handles[wxproto.MsgTypeImage]           = defaultPlainHandle
	handles[wxproto.EventTypeSubscribe]     = defaultEventHandle
	handles[wxproto.EventTypeUnsubscribe]   = defaultEventHandle
}

func (wh *WXSDKHandle)ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("parse form error: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	fmt.Printf("handle http request, method %v\n", r.Method)

    //校验签名
	if !wh.checkSignature(r) {
		fmt.Println("checkSignature error")
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	switch (r.Method) {
	case "GET":
		//处理GET请求
		wh.handleGet(w, r)
	case "POST":
		//处理POST请求
		wh.handlePost(w, r)
	default:
		fmt.Printf("inv http method %v\n", r.Method)
		http.Error(w, "only GET or POST method allowed", http.StatusBadRequest)
	}
}

func (wh *WXSDKHandle) checkSignature(r *http.Request) bool {
	signature := r.FormValue("signature")
	timestamp := r.FormValue("timestamp")
	nonce := r.FormValue("nonce")
	strs := []string{wxproto.Token, timestamp, nonce}
	sort.Strings(strs)
	s := strings.Join(strs, "")

	h := sha1.New()
	h.Write([]byte(s))
	result := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Printf("result: [%v], signature: [%v]\n", result, signature)
	return result == signature
}

func (wh *WXSDKHandle) handleGet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, r.FormValue("echostr"))
}

func (wh *WXSDKHandle) handlePost(w http.ResponseWriter, r *http.Request) {
	var err error
	//解包 r.Body
	var reqMsg wxproto.Message
	err = wh.unpackPkg(r, &reqMsg)
	if err != nil {
		fmt.Printf("unpackPkg err: %v\n", err)
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	fmt.Printf("incoming reqmsg: type %v\n", reqMsg.MsgType)

	var ctx wxproto.RequestContext
	ctx.ToUserName = reqMsg.FromUserName
	ctx.FromUserName = reqMsg.ToUserName
	//ctx.CreateTime = uint32(time.Now().Unix())
	ctx.MsgType = reqMsg.MsgType

	ctx.W = w

	//判断request msgtype
	if reqMsg.MsgType == wxproto.MsgTypeEvent {
		processEvent(&ctx, &reqMsg)
	} else {
		processPlain(&ctx, &reqMsg)
	}
}

func (wh *WXSDKHandle) unpackPkg(r *http.Request, m *wxproto.Message) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, m)
}

func replySuccess(ctx *wxproto.RequestContext) {
	ctx.W.Write([]byte("success"))
}


func replyOK(ctx *wxproto.RequestContext) {
	ctx.W.Write([]byte(""))
}
