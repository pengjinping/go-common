package client

import "C"
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

type Client struct {
	Header      map[string]string      // 头信息
	Data        map[string]interface{} // 数据内容
	Domain      string                 // 域名地址
	IsLog       bool                   // 是否写日志
	Timeout     int                    // 请求超时时间
	DialTimeout int                    // TCP连接的时间
	id          string                 // 请求标识
}

func New() *Client {
	header := make(map[string]string)
	header["Content-Type"] = "application/json"

	return &Client{
		Header:      header,
		Data:        make(map[string]interface{}),
		Domain:      "",
		IsLog:       false,
		Timeout:     10,
		DialTimeout: 3,
		id:          randString(16),
	}
}

func (c *Client) Get(path string) ([]byte, error) {
	url := c.Domain + path

	resp, err := http.Get(url)
	if err != nil {
		return c.writeError(url, fmt.Sprintf("创建请求失败:%v\n", err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.writeError(url, fmt.Sprintf("访问请求失败:%v\n", err))
	}

	c.writeLog(url, body)

	return body, nil
}

func (c *Client) Post(path string) ([]byte, error) {
	var reader io.Reader
	if c.Data != nil {
		data, _ := json.Marshal(&c.Data)
		reader = bytes.NewReader(data)
	}
	url := c.Domain + path

	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return c.writeError(url, fmt.Sprintf("创建请求失败:%v\n", err))
	}

	for key, val := range c.Header {
		request.Header.Set(key, val)
	}

	client := c.makeClient()
	resp, err := client.Do(request)
	if err != nil {
		return c.writeError(url, fmt.Sprintf("访问请求失败:%v\n", err))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.writeError(url, fmt.Sprintf("获取请求数据失败:%v\n", err))
	}

	c.writeLog(url, body)

	return body, nil
}

func (c Client) makeClient() *http.Client {
	timeOut := time.Duration(c.Timeout) * time.Second
	return &http.Client{
		Transport: &http.Transport{
			Dial: func(netWork, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netWork, addr, time.Duration(c.DialTimeout)*time.Second) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				err = conn.SetDeadline(time.Now().Add(timeOut)) //设置发送接受数据超时
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
			ResponseHeaderTimeout: timeOut,
		},
	}
}

func randString(len int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	byteWord := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		byteWord[i] = byte(b)
	}

	return string(byteWord)
}

func (c *Client) writeLog(url string, body []byte) {
	if c.IsLog {
		info := make(map[string]interface{})
		info["url"] = url
		info["params"] = c.Data
		info["response"] = string(body)

		log.Printf("请求信息：%#v", info)
	}
}

func (c *Client) writeError(url string, errMsg string) ([]byte, error) {
	info := make(map[string]interface{})
	info["id"] = c.id
	info["url"] = url
	info["params"] = c.Data
	info["error"] = errMsg

	var out bytes.Buffer
	b, _ := json.Marshal(info)
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		log.Printf("请求失败：\n%#v\n", info)
	} else {
		log.Printf("请求失败：\n%v\n", out.String())
	}

	return nil, fmt.Errorf(errMsg)
}
