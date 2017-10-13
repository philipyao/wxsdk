package wxsdk

import (
    "fmt"
	"time"
    "net/http"
    "image/color"
    "image/png"
    "github.com/afocus/captcha"
    "os"
)

var (
	MenuBtns []Button
    cap *captcha.Captcha

    loginRsp []byte = make([]byte, 1000000)
)

func Init() error {
    f , err := os.Open("./conf/login_conf.json")
    if err != nil{
        fmt.Printf("open err %v\n", err)
        return err
    }
    //文件关闭
    defer f.Close()

    n, err := f.Read(loginRsp)
    if err != nil {
        fmt.Printf("read err %v\n", err)
        return err
    }
    fmt.Printf("read login conf len %v\n", n)
    loginRsp = loginRsp[:n]

    return nil
}

func Serve() {
    var err error

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
    //files
    http.Handle("/", http.FileServer(http.Dir("./dist")))
    fs := http.FileServer(http.Dir("download"))
    http.Handle("/api/download/", http.StripPrefix("/api/download/", fs))

    //验证码
    cap = captcha.New()
    if err = cap.SetFont("comic.ttf"); err != nil {
        panic(err.Error())
    }
    cap.SetSize(128, 64)
    cap.SetDisturbance(captcha.MEDIUM)
    cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
    cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
    http.HandleFunc("/admin/login/captcha/load", func(w http.ResponseWriter, r *http.Request) {
        img, str := cap.Create(6, captcha.ALL)
        png.Encode(w, img)
        fmt.Printf("load captcha, img %v, str %v\n", img, str)
    })

    //后台登录
    handle_admin()

    err = http.ListenAndServe("0.0.0.0:80", nil)
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

