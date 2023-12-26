package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type Img struct {
	ImgLink string
	data    []byte
	Base    string
}

func (i *Img) GetImgLink() error {
	// 发起GET请求获取图片数据
	resp, err := http.Get(i.ImgLink)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// 读取响应体
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	i.data = data
	return nil
}
func (i *Img) EncodeToBase64() {
	i.Base = base64.StdEncoding.EncodeToString(i.data)
}
