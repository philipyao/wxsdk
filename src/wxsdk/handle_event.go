package wxsdk

import (
	"fmt"
    "wxsdk/wxproto"
)

func processEvent(ctx *wxproto.RequestContext, reqMsg *wxproto.Message) {
	event := wxproto.MsgTypeEvent + wxproto.MsgType("." + reqMsg.Event)
	hdl, ok := handles[event]
	if !ok {
		fmt.Printf("unsupported event %v\n", reqMsg.Event)
		replyOK(ctx)
		return
	}
	hdl(ctx, reqMsg)
}

func defaultEventHandle(ctx *wxproto.RequestContext, reqMsg *wxproto.Message) {
	fmt.Printf("[event]%+v\n", reqMsg)
	replyOK(ctx)
}
