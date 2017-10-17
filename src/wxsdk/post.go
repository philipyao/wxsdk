package wxsdk

import (
	"fmt"
    "os"
	"errors"
	"encoding/json"
)

//素材上传返回信息
type TempMediaRsp struct {
	WeiXinRspHeader

	Type        string  `json:"type"`
	MediaId     string  `json:"media_id"`
	CreatedAt   int     `json:"created_at"`
}

//永久素材上传返回
type MediaRsp struct {
	WeiXinRspHeader

	MediaId     string  `json:"media_id"`
	//NewsImgUrl  string  `json:"url,omitempty"`
    Url         string  `json:"url,omitempty"`
}

//图文定义
type Article struct {
	Title               string          `json:"title"`
	ThumbMediaId        string          `json:"thumb_media_id"`     //图文消息的封面图片素材id（必须是永久mediaID）
	Author              string          `json:"author,omitempty"`
	Digest              string          `json:"digest,omitempty"`   //仅有单图文消息才有摘要，多图文此处为空
	ShowCoverPic        int             `json:"show_cover_pic"`     //是否显示封面图片
	Content             string          `json:"content"`            //图文消息的具体内容，支持HTML标签，必须少于2万字符，
                                                                    //小于1M，且此处会去除JS,涉及图片url必须来源"上传图文消息内的图片获取URL"接口获取。
                                                                    //外部图片url将被过滤。
	ContentSourceUrl    string          `json:"content_source_url"` //阅读原文链接
}

type News struct {
	Articles    []*Article      `json:"articles"`
}

// 拉取自定义菜单返回
type MenuRsp struct {
	WeiXinRspHeader

	Menu    *Menu     `json:"menu"`
}

type Menu struct {
	Button []Button `json:"button"`
	MenuId int      `json:"menuid,omitempty"` // 菜单 id
}

// Button 菜单上的按钮
type Button struct {
	Name      string   `json:"name"`
	Type      string   `json:"type,omitempty"`
	Key       string   `json:"key,omitempty"`
	URL       string   `json:"url,omitempty"`
	SubButton []Button `json:"sub_button,omitempty"`
}

// 拉取标签列表
type TagsRsp struct {
    WeiXinRspHeader

    Tags    []*TagEntry         `json:"tags"`
}
type TagEntry struct {
    ID      int         `json:"id"`
    Name    string      `json:"name"`
    Count   int         `json:"count"`
}

// 添加标签
type AddTagReq struct {
    Tag     *TagEntry   `json:"tag"`
}
type AddTagRsp struct {
    WeiXinRspHeader
    Tag     *TagEntry   `json:"tag"`
}

// 批量设置标签
type BatchSetTagReq struct {
    Tagid           int             `json:"tagid"`
    OpenidList      []string        `json:"openid_list"`
}
type BatchSetTagRsp struct {
    WeiXinRspHeader
}

// 拉取粉丝列表
type FansListRsp struct {
    WeiXinRspHeader
    Total       int             `json:"total"`
    Count       int             `json:"count"`
    Data        *FansListData   `json:"data"`
    NextOpenid  string          `json:"next_openid"`
}
type FansListData struct {
    Openid      []string        `json:"openid"`
}

// svr <==> wxsvr
type ListTagFansReq struct {
    Tagid       int             `json:"tagid"`
    NextOpenid  string          `json:"next_openid"`
}
type ListTagFansRsp struct {
    WeiXinRspHeader
    Count       int             `json:"count"`
    Data        *FansListData   `json:"data"`
    NextOpenid  string          `json:"next_openid"`
}

//批量拉取粉丝信息
type BatchGetFansInfoReq struct {
    UserList   []*FansOpenid       `json:"user_list"`
}
type FansOpenid struct {
    Openid      string          `json:"openid"`
}
type BatchGetFansInfoRsp struct {
    WeiXinRspHeader
    UserInfoList    []*FansInfo     `json:"user_info_list"`
}
type FansInfo struct {
    Subscribe           int         `json:"subscribe"`
    Openid              string      `json:"openid"`
    Unionid             string      `json:"unionid"`
    Nickname            string      `json:"nickname"`
    Sex                 int         `json:"sex"`
    Language            string      `json:"language"`
    City                string      `json:"city"`
    Province            string      `json:"province"`
    Country             string      `json:"country"`
    Headimgurl          string      `json:"headimgurl"`
    SubscribeTime       int         `json:"subscribe_time"`
    Remark              string      `json:"remark"`     //备注
    Groupid             int         `json:"groupid"`
    TagidList           []int       `json:"tagid_list"`
}

type WXListMaterialReq struct {
    Type                string      `json:"type"`
    Offset              int         `json:"offset"`
    Count               int         `json:"count"`
}
type WXListMaterialRsp struct {
    WeiXinRspHeader
    TotalCount          int                     `json:"total_count"`
    ItemCount           int                     `json:"item_count"`
    Item                []*MaterialProfile      `json:"item"`
}
type Materials struct {
    TotalCount          int
    ItemCount           int
    Item                []*MaterialProfile
}

//下载特定的素材
type WXGetMaterialReq struct {
    MediaID             string              `json:"media_id"`
}

//////////////////////////////////////////////////////////////

//////////  菜单管理  ////////////////////////

// CreateMenu 创建自定义菜单
func CreateMenu(buttons []Button) error {
	if len(buttons) > 3 {
		return errors.New("too many first level menu, must less than 3")
	}
	for _, sub := range buttons {
		if len(sub.SubButton) > 5 {
			return errors.New("too many second level menu, must less than 5")
		}
	}

	menu := struct {
		Button []Button `json:"button"`
	}{buttons}

	url := fmt.Sprintf(UrlMenuCreate, AccessToken())
	resp, err := postJson(url, menu)
	if err != nil {
		return err
	}

	var mrsp MenuRsp
	err = json.Unmarshal(resp, &mrsp)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if mrsp.ErrorCode != 0 {
		err = fmt.Errorf("postMenu error %v, %v", mrsp.ErrorCode, mrsp.ErrorMessage)
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("update menu success.")
	return nil
}

// getMenu查询菜单
func GetMenu() (*Menu, error) {
	url := fmt.Sprintf(UrlMenuGet, AccessToken())

	var mrsp MenuRsp
	err := getJson(url, &mrsp)
	if err != nil {
		return nil, err
	}

	if mrsp.ErrorCode != 0 {
		err = fmt.Errorf("getMenu error %v, %v", mrsp.ErrorCode, mrsp.ErrorMessage)
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("GetMenu success.")
	return mrsp.Menu, nil
}

// DeleteMenu 删除菜单
func DeleteMenu() (err error) {
	url := fmt.Sprintf(UrlMenuDelete, AccessToken())
	var mrsp MenuRsp
	err = getJson(url, &mrsp)
	if err != nil {
		return err
	}

	if mrsp.ErrorCode != 0 {
		err = fmt.Errorf("getMenu error %v, %v", mrsp.ErrorCode, mrsp.ErrorMessage)
		fmt.Println(err.Error())
		return err
	}
	fmt.Println("DeleteMenu success.")
	return nil
}


//////////  用户标签  ////////////////////////

func ListTag(rsp *TagsRsp) (err error) {
    url := fmt.Sprintf(UrlTagList, AccessToken())
    err = getJson(url, rsp)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("ListTag error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return err
    }
    fmt.Printf("ListTag success. %+v\n", rsp.Tags)
    return nil
}

func AddTag(name string, rsp *AddTagRsp) (err error) {
    url := fmt.Sprintf(UrlTagAdd, AccessToken())
    var req AddTagReq
    req.Tag = &TagEntry{Name: name}
    data, err := postJson(url, req)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    err = json.Unmarshal(data, rsp)
    if err != nil {
        fmt.Println(err)
        return err
    }

    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("AddTag error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return err
    }
    fmt.Printf("AddTag success. %+v\n", rsp.Tag)
    return nil
}

func BatchSetTag(req *AdminTagBatchSetReq, rsp *AdminTagBatchSetRsp) (err error) {
    url := fmt.Sprintf(UrlTagBatchSet, AccessToken())

    var onereq BatchSetTagReq
    onereq.OpenidList = req.OpenidList
    for _, tagid := range req.Tagids {
        onereq.Tagid = tagid
        data, err := postJson(url, onereq)
        if err != nil {
            fmt.Println(err.Error())
            return err
        }

        var onersp BatchSetTagRsp
        err = json.Unmarshal(data, &onersp)
        if err != nil {
            fmt.Println(err)
            return err
        }

        if onersp.ErrorCode != 0 {
            err = fmt.Errorf("BatchSetTag<%v> error %v, %v", tagid, onersp.ErrorCode, onersp.ErrorMessage)
            fmt.Println(err.Error())
            rsp.ErrCode = onersp.ErrorCode
            return err
        }
        fmt.Printf("BatchSetTag<%v> done\n", tagid)
    }

    fmt.Println("BatchSetTag success")
    return nil
}

//////////  粉丝  /////////////////////////////

//==== 全部粉丝列表
func ListFans(rsp *FansListRsp) (err error) {
    url := fmt.Sprintf(UrlFansList, AccessToken())
    err = getJson(url, rsp)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("ListTag error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return err
    }
    fmt.Println("ListFans success. %+v", rsp)
    return nil
}

//==== 特定tag粉丝列表
func ListTagFans(tagid int, rsp *ListTagFansRsp) (err error) {
    url := fmt.Sprintf(UrlTagFansList, AccessToken())

    var req ListTagFansReq
    req.Tagid = tagid

    data, err := postJson(url, &req)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    err = json.Unmarshal(data, rsp)
    if err != nil {
        fmt.Println(err)
        return err
    }
    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("ListTagFans error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return err
    }

    fmt.Println("ListTagFans success.")
    return nil
}

//==== 批量
func BatchGetFansInfo(openids []string, rsp *BatchGetFansInfoRsp) (err error) {
    url := fmt.Sprintf(UrlFansBatchGet, AccessToken())
    var req BatchGetFansInfoReq
    req.UserList = make([]*FansOpenid, len(openids))
    for i, openid := range openids {
        req.UserList[i] = &FansOpenid{Openid: openid}
    }
    data, err := postJson(url, req)
    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    err = json.Unmarshal(data, rsp)
    if err != nil {
        fmt.Println(err)
        return err
    }
    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("BatchGetFansData error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return err
    }

    fmt.Println("BatchGetFansData success.")
    return nil
}


//////////  素材管理  ////////////////////////

// 拉取普通永久素材列表（图片、语音、视频）
func ListMaterial(mtype string, offset, count int) (materials *Materials, err error) {
    fmt.Printf("ListMaterial: mtype %v, offset %v, count %v\n", mtype, offset, count)
    url := fmt.Sprintf(UrlMediaList, AccessToken())
    var req WXListMaterialReq
    req.Type = mtype
    req.Offset = offset
    req.Count = count
    data, err := postJson(url, req)
    if err != nil {
        fmt.Println(err.Error())
        return nil, err
    }

    var rsp WXListMaterialRsp
    err = json.Unmarshal(data, &rsp)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    if rsp.ErrorCode != 0 {
        err = fmt.Errorf("ListMaterial error %v, %v", rsp.ErrorCode, rsp.ErrorMessage)
        fmt.Println(err.Error())
        return nil, err
    }

    fmt.Printf("ListMaterial success. total %v, count %v\n", rsp.TotalCount, rsp.ItemCount)
    return &Materials{
        TotalCount: rsp.TotalCount,
        ItemCount: rsp.ItemCount,
        Item: rsp.Item,
    }, nil
}

func GetMaterial(mediaID string) (string, error) {
    fmt.Printf("GetMaterial: mediaID %v\n", mediaID)
    url := fmt.Sprintf(UrlMediaGet, AccessToken())
    var req WXGetMaterialReq
    req.MediaID = mediaID
    content, contentType, err := requestWeiXin(url, req)
    if err != nil {
        fmt.Println(err.Error())
        return "", err
    }

    if contentType == ContentTypeText {
        //有文本json格式, 错误返回
        errMsg := &WeiXinRspHeader{}
        err = json.Unmarshal(content, errMsg)
    } else {
        //直接是媒体数据流，下载到服务器本地
        fmt.Printf("write media(%v) to file, length %v\n", mediaID, len(content))
        f, err := os.OpenFile(fmt.Sprintf("./download/%v", mediaID), os.O_WRONLY | os.O_CREATE , 0644)
        if err != nil {
            return "", err
        }
        defer f.Close()
        n, err := f.Write(content)
        if n != len(content) {
            fmt.Println("write error")
            return "", fmt.Errorf("write length error %d %d", n, len(content))
        }
        f.Sync()
    }
    return mediaID, nil
}

//上传临时素材
//返回微信生成的media_id
func postTempMedia(filename, mediaType string) (string, error) {
    targetUrl := fmt.Sprintf(UrlTempMediaUpload, AccessToken(), mediaType)
    resp, err := postFile(targetUrl, filename, "media", nil)
    if err != nil {
        return "", err
    }

    var mediaRsp TempMediaRsp
    err = json.Unmarshal(resp, &mediaRsp)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    if mediaRsp.ErrorCode != 0 {
        err = fmt.Errorf("postTempMedia error %v, %v", mediaRsp.ErrorCode, mediaRsp.ErrorMessage)
        fmt.Println(err.Error())
        return "", err
    }
    fmt.Printf("postTempMedia success: filename %v, mediaid %v\n", filename, mediaRsp.MediaId)
    return mediaRsp.MediaId, nil
}


//上传永久普通素材(图片、语音)
//返回微信生成的media_id
func postMedia(filename, mediaType string) (string, error) {
    targetUrl := fmt.Sprintf(UrlMediaUpload, AccessToken(), mediaType)
    resp, err := postFile(targetUrl, filename, "media", nil)
    if err != nil {
        return "", err
    }

    var mediaRsp MediaRsp
    err = json.Unmarshal(resp, &mediaRsp)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    if mediaRsp.ErrorCode != 0 {
        err = fmt.Errorf("postMedia error %v, %v", mediaRsp.ErrorCode, mediaRsp.ErrorMessage)
        fmt.Println(err.Error())
        return "", err
    }
    fmt.Printf("postMedia success: filename %v, mediaid %v\n", filename, mediaRsp.MediaId)
    return mediaRsp.MediaId, nil
}

//上传永久视频素材
//返回微信生成的media_id
func postVideo(filename, title, introduction string) (string, error) {
    targetUrl := fmt.Sprintf(UrlMediaUpload, AccessToken(), MediaTypeVideo)
    params := map[string]string{
        "title": title,
        "introduction": introduction,
    }
    resp, err := postFile(targetUrl, filename, "media", params)
    if err != nil {
        return "", err
    }

    var mediaRsp MediaRsp
    err = json.Unmarshal(resp, &mediaRsp)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    if mediaRsp.ErrorCode != 0 {
        err = fmt.Errorf("postVideo error %v, %v", mediaRsp.ErrorCode, mediaRsp.ErrorMessage)
        fmt.Println(err.Error())
        return "", err
    }
    fmt.Printf("postVideo success: filename %v, mediaid %v\n", filename, mediaRsp.MediaId)
    return mediaRsp.MediaId, nil
}

//上传图片（永久图文素材内的图片）素材
//不占用公众号的素材库中图片数量的5000个的限制
//返回微信生成的图片url
func postNewsImg(filename string) (string, error) {
    targetUrl := fmt.Sprintf(UrlMediaUploadNewsImg, AccessToken())
    resp, err := postFile(targetUrl, filename, "media", nil)
    if err != nil {
        return "", err
    }

    var mediaRsp MediaRsp
    err = json.Unmarshal(resp, &mediaRsp)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    if mediaRsp.ErrorCode != 0 {
        err = fmt.Errorf("postNewsImg error %v, %v", mediaRsp.ErrorCode, mediaRsp.ErrorMessage)
        fmt.Println(err.Error())
        return "", err
    }
    fmt.Printf("postNewsImg success: filename %v, NewsImgUrl %v\n", filename, mediaRsp.Url)
    return mediaRsp.Url, nil
}


//上传永久图文素材
//返回微信生成的media_id
func postNews(news *News) (string, error) {
    targetUrl := fmt.Sprintf(UrlMediaUploadNews, AccessToken())
    reply, err := postJson(targetUrl, news)
    if err != nil {
        return "", err
    }
    var mediaRsp MediaRsp
    err = json.Unmarshal(reply, &mediaRsp)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    if mediaRsp.ErrorCode != 0 {
        err = fmt.Errorf("postVideo error %v, %v", mediaRsp.ErrorCode, mediaRsp.ErrorMessage)
        fmt.Println(err.Error())
        return "", err
    }
    fmt.Printf("postNews success: mediaid %v\n", mediaRsp.MediaId)
    return mediaRsp.MediaId, nil
}




//////////////////////////////////////////////////////////

