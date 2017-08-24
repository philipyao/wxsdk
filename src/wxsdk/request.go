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
	resp, err := http.Post(targetUrl, "application/json", reader)
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

//上传文件，支持附带额外参数
func postFile(targetUrl string, filename string, params map[string]string) ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filename)
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

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
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

