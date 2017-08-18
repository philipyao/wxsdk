package wxsdk

import(
    "net/http"
)

type WXSDKHandle struct {}

var defaultServeMux WXSDKHandle

func (wh *WXSDKHandle)ServeHTTP(http.ResponseWriter, *http.Request) {
    //校验签名

    //处理GET请求

    //处理POST请求
}
