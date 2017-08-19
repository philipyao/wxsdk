package wxsdk

import (
	"fmt"
	"time"
	//"encoding/xml"
)

func processPlain(ctx *RequestContext, reqMsg *Message) {
	hdl, ok := handles[reqMsg.MsgType]
	if !ok {

	}
	hdl(ctx, reqMsg)
}


func defaultTextHandle(ctx *RequestContext, reqMsg *Message) {
	fmt.Printf("[text]%+v\n", reqMsg)
	replySuccess(ctx)
}

func defaultImageHandle(ctx *RequestContext, reqMsg *Message) {
	fmt.Printf("[image]%+v\n", reqMsg)
	replySuccess(ctx)
}


func replySuccess(ctx *RequestContext) {
	ctx.w.Write([]byte("success"))
}

func replyText(ctx *RequestContext, content string) {
	doReply(ctx, buildText(content))
}

func doReply(ctx *RequestContext, rspMsg *Message) {
	rspMsg.CreateTime = uint32(time.Now().Unix())
	rspMsg.ToUserName = ctx.ToUserName
	rspMsg.FromUserName = ctx.FromUserName
}

func buildText(content string) *Message {
	var rspMsg Message
	rspMsg.MsgType = MsgTypeText
	rspMsg.Content = content
	return &rspMsg
}