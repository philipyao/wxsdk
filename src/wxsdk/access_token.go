package wxsdk

import(
    "sync"
)

//获取AccessToken
var AccessToken func() string

//维护AccessToken
func keepAccessToken() {
    var accessToken struct{
        token   string
        expireIn int
        lock    *sync.RWMutex
    }

    AccessToken = func() string {
        accessToken.lock.RLock()
        defer accessToken.lock.RUnlock()
        return accessToken.token
    }

    //TODO for 循环刷新
}
