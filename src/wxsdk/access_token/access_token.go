package access_token

import(
    "fmt"
    "sync"
    "time"

    "wxsdk/utils"
)

const (
    //测试号
    //APPID           = "wx4d953a6ad5112405"
    //APPSecret       = "56657d6773e91ee3a2814273b9813fc0"

    APPID           = "wxe5c18b742279f207"
    APPSecret       = "9153504404a33eb02e3b2298b58111fd"
    UrlAccessToken  = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v"
)

type accessToken struct{
    token   string
    sync.RWMutex
}

type TokenRsp struct {
    ErrorCode    int            `json:"errcode,omitempty"`
    ErrorMessage string         `json:"errmsg,omitempty"`

    Token   string              `json:"access_token"`
    ExpireIn int                `json:"expires_in"`
}

// =========对外接口==========

//获取AccessToken
var AccessToken func() string
//主动刷新token
var RefreshAccessToken func() (int, error)


//========================================================
//维护AccessToken
func StartAccessToken() {
    var _token = new(accessToken)

    AccessToken = func() string {
        _token.RLock()
        defer _token.RUnlock()
        return _token.token
    }

    RefreshAccessToken = func() (int, error) {
        token, expirein, err := requestAccessToken()
        if err != nil {
            fmt.Printf("requestAccessToken error %v\n", err)
            return 0, err
        }
        fmt.Printf("requestAccessToken success: expire %v, token %v\n", expirein, token)
        _token.Lock()
        _token.token = token
        _token.Unlock()

        return expirein, nil
    }

    // for 循环刷新
    go func() {
        for {
            expirein, err := RefreshAccessToken()
            if err != nil {
                time.Sleep(time.Minute)
                continue
            }

            timer := time.NewTimer(time.Second * time.Duration(expirein))
            <- timer.C
        }
    }()
}

func requestAccessToken() (string, int, error) {
    url := fmt.Sprintf(UrlAccessToken, APPID, APPSecret)
    var rsp TokenRsp
    err := utils.GetJson(url, &rsp)
    if err != nil {
        return "", 0, err
    }
    if rsp.ErrorCode != 0 {
        return "", 0, fmt.Errorf("requestAccessToken error: %v %v", rsp.ErrorCode, rsp.ErrorMessage)
    }
    return rsp.Token, rsp.ExpireIn, nil
}

