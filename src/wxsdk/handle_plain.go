package wxsdk

import (
	"fmt"
	"time"
	"encoding/xml"
    "wxsdk/wxproto"
)

func processPlain(ctx *wxproto.RequestContext, reqMsg *wxproto.Message) {
	hdl, ok := handles[reqMsg.MsgType]
	if !ok {
		fmt.Printf("unsupported msgtype %v\n", reqMsg.MsgType)
		replySuccess(ctx)
		return
	}
	hdl(ctx, reqMsg)
}


func defaultPlainHandle(ctx *wxproto.RequestContext, reqMsg *wxproto.Message) {
	fmt.Printf("[plain]%+v\n", reqMsg)
	replyOK(ctx)
}

func replyText(ctx *wxproto.RequestContext, content string) {
	doReply(ctx, buildText(content))
}

func doReply(ctx *wxproto.RequestContext, rspMsg *wxproto.Message) {
	rspMsg.CreateTime = uint32(time.Now().Unix())
	rspMsg.ToUserName = ctx.ToUserName
	rspMsg.FromUserName = ctx.FromUserName

	data, err := xml.MarshalIndent(rspMsg, "", "    ")
	if err != nil {
		fmt.Printf("marshal error %v\n", err)
		replySuccess(ctx)
		return
	}
	//reply消息给微信服务器
	ctx.W.Header().Set("Content-Type", "text/xml; charset=utf-8")
	ctx.W.Write(data)
}

func buildText(content string) *wxproto.Message {
	var rspMsg wxproto.Message
	rspMsg.MsgType = wxproto.MsgTypeText
	rspMsg.Content = content
	return &rspMsg
}