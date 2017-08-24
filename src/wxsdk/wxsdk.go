package wxsdk

import (
    "fmt"
	"time"
    "net/http"
)

var (
	MenuBtns []Button
)

func Init() error {
    return nil
}

func Serve() {
    // 维护AccessToken
    keepAccessToken()

	go func() {
		time.Sleep(5 * time.Second)
		if len(MenuBtns) > 0 {
			err := CreateMenu(MenuBtns)
			if err != nil {
				fmt.Printf("CustomsizeMenu failed: %v\n", err)
			}
		}
	}()

    //http.Serve
    http.Handle("/weixin", &defaultServeMux)
    err := http.ListenAndServe("0.0.0.0:80", nil)
    if err != nil {
        fmt.Printf("ListenAndServe error: %v\n", err)
    }
}

func CustomsizeMenu(buttons []Button) {
	MenuBtns = buttons
}

//自定义普通消息回复
func HandlePlain(hdl interface{}) {
}

//自定义事件处理
func HandleEvent(hdl interface{}) {
}

//自定义推送
func HandlePush(hdl interface{}) {
}

