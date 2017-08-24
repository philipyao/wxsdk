package wxsdk

import (
	"fmt"
)

func processEvent(ctx *RequestContext, reqMsg *Message) {
	event := MsgTypeEvent + MsgType("." + reqMsg.Event)
	hdl, ok := handles[event]
	if !ok {
		fmt.Printf("unsupported event %v\n", reqMsg.Event)
		replyOK(ctx)
		return
	}
	hdl(ctx, reqMsg)
}

func defaultEventHandle(ctx *RequestContext, reqMsg *Message) {
	fmt.Printf("[event]%+v\n", reqMsg)
	replyOK(ctx)
}
