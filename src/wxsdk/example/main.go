package main

import (
    "fmt"
    "wxsdk"
)

func main() {
    fmt.Println("start wxsdk service.")
    wxsdk.Init()
	ExampleCreateMenu()
    wxsdk.Serve()
}


func ExampleCreateMenu() {
	buttons := []wxsdk.Button{
		wxsdk.Button{
			Name: "扫码",
			SubButton: []wxsdk.Button{
				wxsdk.Button{
					Name: "扫码带提示",
					Type: wxsdk.MenuTypeScancodeWaitmsg,
					Key:  "rselfmenu_0_0",
				},
				wxsdk.Button{
					Name: "扫码推事件",
					Type: wxsdk.MenuTypeScancodePush,
					Key:  "rselfmenu_0_1",
				},
			},
		},
		wxsdk.Button{
			Name: "发图",
			SubButton: []wxsdk.Button{
				wxsdk.Button{
					Name: "系统拍照发图",
					Type: wxsdk.MenuTypePicSysphoto,
					Key:  "rselfmenu_1_0",
				},
				wxsdk.Button{
					Name: "拍照或者相册发图",
					Type: wxsdk.MenuTypePicPhotoOrAlbum,
					Key:  "rselfmenu_1_1",
				},
				wxsdk.Button{
					Name: "微信相册发图",
					Type: wxsdk.MenuTypePicWeixin,
					Key:  "rselfmenu_1_2",
				},
			},
		},
		wxsdk.Button{
			Name: "测试",
			SubButton: []wxsdk.Button{
				wxsdk.Button{
					Name: "腾讯",
					Type: wxsdk.MenuTypeView,
					URL:  "http://qq.com",
				},
				wxsdk.Button{
					Name: "发送位置",
					Type: wxsdk.MenuTypeLocationSelect,
					Key:  "rselfmenu_2_0",
				},
			},
		},
	}

	wxsdk.CustomsizeMenu(buttons)
}

