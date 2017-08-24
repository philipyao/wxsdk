package wxsdk

import (
	"fmt"
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
	NewsImgUrl  string  `json:"url,omitempty"`
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
	MenuId int      `json:"menuid"` // 菜单 id
}

// Button 菜单上的按钮
type Button struct {
	Name      string   `json:"name"`
	Type      string   `json:"type,omitempty"`
	Key       string   `json:"key,omitempty"`
	URL       string   `json:"url,omitempty"`
	SubButton []Button `json:"sub_button,omitempty"`
}


//上传临时素材
//返回微信生成的media_id
func postTempMedia(filename, mediaType string) (string, error) {
	targetUrl := fmt.Sprintf(UrlTempMediaUpload, AccessToken(), mediaType)
	resp, err := postFile(targetUrl, filename, nil)
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


//上传永久普通素材
//返回微信生成的media_id
func postMedia(filename, mediaType string) (string, error) {
	targetUrl := fmt.Sprintf(UrlMediaUpload, AccessToken(), mediaType)
	resp, err := postFile(targetUrl, filename, nil)
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
	resp, err := postFile(targetUrl, filename, params)
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
	resp, err := postFile(targetUrl, filename, nil)
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
	fmt.Printf("postNewsImg success: filename %v, NewsImgUrl %v\n", filename, mediaRsp.NewsImgUrl)
	return mediaRsp.NewsImgUrl, nil
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
	fmt.Println("postMenu success.")
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
	fmt.Println("getMenu success: %+v", mrsp.Menu)
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


//////////////////////////////////////////////////////////

