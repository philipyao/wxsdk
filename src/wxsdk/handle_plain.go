package wxsdk

import (
	"fmt"
	"time"
	"encoding/xml"
)

func processPlain(ctx *RequestContext, reqMsg *Message) {
	hdl, ok := handles[reqMsg.MsgType]
	if !ok {
		fmt.Printf("unsupported msgtype %v\n", reqMsg.MsgType)
		replySuccess(ctx)
		return
	}
	hdl(ctx, reqMsg)
}


func defaultPlainHandle(ctx *RequestContext, reqMsg *Message) {
	fmt.Printf("[plain]%+v\n", reqMsg)
	replyOK(ctx)
}

func replyText(ctx *RequestContext, content string) {
	doReply(ctx, buildText(content))
}

func doReply(ctx *RequestContext, rspMsg *Message) {
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
	ctx.w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	ctx.w.Write(data)
}

func buildText(content string) *Message {
	var rspMsg Message
	rspMsg.MsgType = MsgTypeText
	rspMsg.Content = content
	return &rspMsg
}