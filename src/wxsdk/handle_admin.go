package wxsdk

import(
    "fmt"
    "net/http"
    "encoding/json"
    "strconv"
    "io/ioutil"
    "io"
    "os"
    "wxsdk/wxproto"
)

const (
    DefaultListMaterialCount    = 20
)

type AdminFansListRsp struct {
    Errcode         int             `json:"errcode"`
    Total           int             `json:"total"`
    Count           int             `json:"count"`
    NextOpenid      string          `json:"next_openid"`
    Fans            []*wxproto.FansInfo     `json:"fans"`
}

type AdminTagBatchSetReq struct {
    Tagids          []int           `json:"tagids"`
    OpenidList      []string        `json:"openid_list"`
}
type AdminTagBatchSetRsp struct {
    ErrCode    int    `json:"errcode"`
    ErrMessage string `json:"errmsg"`
}

//相应的请求为表单提交
type AdminMaterialCreateRsp struct {
    MediaID         string          `json:"media_id"`
}

//相应的请求为表单提交
type AdminMaterialListRsp struct {
    ErrCode         int    `json:"errcode"`
    ErrMessage      string `json:"errmsg"`
    Materials       []*wxproto.MaterialProfile  `json:"materials"`
}

//图片、音频、视频都返回url来提供预览
type AdminMaterialGetRsp struct {
    ErrCode         int             `json:"errcode"`
    ErrMessage      string          `json:"errmsg"`
    URL             string          `json:"url"`
}

type AdminMenuUpdateReq struct {
    Button          []wxproto.Button        `json:"button"`
}

type AdminMenuUpdateRsp struct {
    ErrCode         int             `json:"errcode"`
    ErrMessage      string          `json:"errmsg"`
}

func handle_admin() {

    http.HandleFunc("/api/admin/login", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        userName, passwd, veriCode := r.FormValue("username"), r.FormValue("password"), r.FormValue("code")
        fmt.Printf("ADMIN LOGIN: [%v] [%v] [%v]\n", userName, passwd, veriCode)
        w.Header().Set("Content-Type", "application/json")
        w.Write(loginRsp)
    })

    http.HandleFunc("/api/admin/wechat/tag/list", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        var rsp wxproto.TagsRsp
        err := wxproto.ListTag(&rsp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        doWriteJson(w, rsp)
    })
    http.HandleFunc("/api/admin/wechat/tag/add", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        tagname := r.FormValue("name")
        var rsp wxproto.AddTagRsp
        err = wxproto.AddTag(tagname, &rsp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        doWriteJson(w, rsp)
    })
    http.HandleFunc("/api/admin/wechat/tag/batch_set", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
        //body不用表单，数据需要自己从body解析
        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Printf("read body error %v\n", err)
            return
        }
        if len(reqdata) == 0 {
            http.Error(w, "no reqdata for tag/batch_set", http.StatusBadRequest)
            return
        }
        var req AdminTagBatchSetReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            http.Error(w, "error parse json reqdata for tag/batch_set", http.StatusBadRequest)
            return
        }

        var wxreq wxproto.WXTagBatchSetReq
        wxreq.Tagids = req.Tagids
        wxreq.OpenidList = req.OpenidList
        var wxrsp wxproto.WeiXinRspHeader
        err = wxproto.BatchSetTag(&wxreq, &wxrsp)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        doWriteJson(w, wxrsp)
    })

    http.HandleFunc("/api/admin/wechat/fans/list", func(w http.ResponseWriter, r *http.Request) {
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }

        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        tagidstr := r.FormValue("tagid")
        tagid, err := strconv.Atoi(tagidstr)
        if err != nil {
            http.Error(w, "inv tagid", http.StatusBadRequest)
            return
        }
        fmt.Printf("fans list: tagid<%v>\n", tagid)
        if tagid == 0 {
            //拉取所有粉丝信息
            var lrsp wxproto.FansListRsp
            err := wxproto.ListFans(&lrsp)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            var rsp AdminFansListRsp
            rsp.Errcode = 0
            rsp.Total = lrsp.Total
            rsp.Count = lrsp.Count
            rsp.NextOpenid = lrsp.NextOpenid

            if len(lrsp.Data.Openid) > 0 {
                var irsp wxproto.BatchGetFansInfoRsp
                err = wxproto.BatchGetFansInfo(lrsp.Data.Openid, &irsp)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusBadRequest)
                    return
                }
                rsp.Fans = irsp.UserInfoList
            }

            doWriteJson(w, rsp)
        } else {
            //单独拉取某个tag的粉丝信息
            var lrsp wxproto.ListTagFansRsp
            err := wxproto.ListTagFans(tagid, &lrsp)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            var rsp AdminFansListRsp
            rsp.Errcode = 0
            rsp.Total = 0
            rsp.Count = lrsp.Count
            rsp.NextOpenid = lrsp.NextOpenid

            if len(lrsp.Data.Openid) > 0 {
                var irsp wxproto.BatchGetFansInfoRsp
                err = wxproto.BatchGetFansInfo(lrsp.Data.Openid, &irsp)
                if err != nil {
                    http.Error(w, err.Error(), http.StatusBadRequest)
                    return
                }
                rsp.Fans = irsp.UserInfoList
            }

            doWriteJson(w, rsp)
        }
    })

    //处理用户上传
    http.HandleFunc("/api/admin/upload", func(w http.ResponseWriter, r *http.Request) {
        //内存存储32M，其他放临时文件中
        r.ParseMultipartForm(32 << 20)
        //iview upload组件名字默认为"file"
        file, handler, err := r.FormFile("file")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()
        fmt.Printf("recv upload file: %v\n", handler.Filename)
        f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)  // 此处假设当前目录下已存在upload目录
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
    })

    http.HandleFunc("/api/admin/wechat/material/create", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }

        mtype := r.FormValue("mtype")
        if mtype != "image" && mtype != "voice" && mtype != "video" && mtype != "thumb" {
            errmsg := fmt.Sprintf("unsupported material type: %v\n", mtype)
            fmt.Println(errmsg)
            http.Error(w, errmsg, http.StatusBadRequest)
            return
        }
        mname := r.FormValue("mname")
        mediaID, err := wxproto.PostMedia("./upload/" + mname, mtype)
        if err != nil {
            fmt.Printf("postMedia error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }
        var rsp AdminMaterialCreateRsp
        rsp.MediaID = mediaID
        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/admin/wechat/material/list", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }

        mtype := r.FormValue("mtype")
        if mtype != "image" && mtype != "voice" && mtype != "video" {
            errmsg := fmt.Sprintf("unsupported material type: %v\n", mtype)
            fmt.Println(errmsg)
            http.Error(w, errmsg, http.StatusBadRequest)
            return
        }

        var rsp AdminMaterialListRsp

        m, err := wxproto.ListMaterial(mtype, 0, DefaultListMaterialCount)
        if err != nil {
            fmt.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        rsp.Materials = append(rsp.Materials, m.Item...)
        for len(rsp.Materials) < m.TotalCount {
            tmpm, err := wxproto.ListMaterial(mtype, len(rsp.Materials), DefaultListMaterialCount)
            if err != nil {
                fmt.Println(err)
                http.Error(w, err.Error(), http.StatusBadRequest)
                break
            }
            rsp.Materials = append(rsp.Materials, tmpm.Item...)
        }

        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/admin/wechat/material/get", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
        err := r.ParseForm()
        if err != nil {
            fmt.Printf("parse form error: %v\n", err)
            http.Error(w, "", http.StatusBadRequest)
            return
        }

        mediaID := r.FormValue("media_id")

        var rsp AdminMaterialGetRsp

        url, err := wxproto.GetMaterial(mediaID)
        if err != nil {
            fmt.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        rsp.URL = url

        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/admin/wechat/menu/get", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "GET" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        menu, err := wxproto.GetMenu()
        if err != nil {
            fmt.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        doWriteJson(w, menu)
    })

    http.HandleFunc("/api/admin/wechat/menu/update", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }
        //body不用表单，数据需要自己从body解析
        reqdata, err := ioutil.ReadAll(r.Body)
        if err != nil {
            fmt.Printf("read body error %v\n", err)
            return
        }
        if len(reqdata) == 0 {
            http.Error(w, "no reqdata for menu/update", http.StatusBadRequest)
            return
        }
        var req AdminMenuUpdateReq
        err = json.Unmarshal(reqdata, &req)
        if err != nil {
            http.Error(w, "error parse json reqdata for tag/batch_set", http.StatusBadRequest)
            return
        }

        var rsp AdminMenuUpdateRsp
        err = wxproto.CreateMenu(req.Button)
        if err != nil {
            rsp.ErrCode = -1
            rsp.ErrMessage = err.Error()
        }
        doWriteJson(w, rsp)
    })

    http.HandleFunc("/api/admin/wechat/menu/delete", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            fmt.Printf("handle http request, method %v\n", r.Method)
            http.Error(w, "inv method", http.StatusBadRequest)
            return
        }

        err := wxproto.DeleteMenu()
        if err != nil {
            fmt.Println(err)
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
    })
}


func doWriteJson(w http.ResponseWriter, pkg interface{}) {
    data, err := json.Marshal(pkg)
    if err != nil {
        fmt.Println(err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}

func doWriteError(w http.ResponseWriter, errmsg string) {
    w.Write([]byte(errmsg))
}
