package wxsdk

import (
    "fmt"
    "net/http"
)

func Init() error {
    return nil
}

func Serve() {
    // 维护AccessToken
    keepAccessToken()

    //http.Serve
    http.Handle("/weixin", &defaultServeMux)
    err := http.ListenAndServe("0.0.0.0:80", nil)
    if err != nil {
        fmt.Printf("ListenAndServe error: %v\n", err)
    }
}

//自定义消息回复
func HandleMsg(hdl interface{}) {
}

//自定义事件处理
func HandleEvent(hdl interface{}) {
}

//自定义推送
func HandlePush(hdl interface{}) {
}

