package wxsdk

import (
	"os"
	"fmt"
	"io"
	"io/ioutil"
	"bytes"
	"net/http"
	"mime/multipart"
	"encoding/json"
)

const (
    ContentTypeJson         = "application/json"
    ContentTypeText         = "text/plain"
)

func getJson(targetUrl string, reply interface{}) error {
	resp, err := http.Get(targetUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("inv http response status %v", resp.Status)
	}

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Printf("getJson: url %v, rsplen %v\n", targetUrl, len(resp_body))
	err = json.Unmarshal(resp_body, reply)
	if err != nil {
		return err
	}
	return nil
}


func postJson(targetUrl string, pkg interface{}) ([]byte, error) {
	data, err := json.Marshal(pkg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	reader := bytes.NewBuffer(data)
	resp, err := http.Post(targetUrl, ContentTypeJson, reader)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.Status)
	}

	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return resp_body, nil
}

//给WX服务器发送请求，请求包一般为json格式
//返回格式为文本格式或者二进制，需要区别处理
func requestWeiXin(url string, pkg interface{}) (content []byte, contentType string, err error) {
    data, err := json.Marshal(pkg)
    if err != nil {
        fmt.Println(err)
        return nil, "", err
    }
    reader := bytes.NewBuffer(data)
    resp, err := http.Post(url, ContentTypeJson, reader)
    if err != nil {
        return nil, "", err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, "", fmt.Errorf("%v", resp.Status)
    }

    defer resp.Body.Close()
    resp_body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, "", err
    }
    return resp_body, resp.Header.Get("Content-Type"), nil
}

//上传文件，支持附带额外参数
func postFile(url, filename, fieldname string, params map[string]string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(fieldname, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil, err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil, err
	}
	defer fh.Close()

	//写入文件数据
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil, err
	}
	//写入参数
	if params != nil {
		for key, val := range params {
			bodyWriter.WriteField(key, val)
		}
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v", resp.Status)
	}

	return resp_body, nil
}

