package admin_handler

import (
	"encoding/json"
	"fmt"
	"github.com/henrylee2cn/faygo"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

/*
封装各种公共的方法,直接在使用的结构体里面嵌套base使用即可
*/

type data_info map[string]interface{}

type jsonData struct {
	Code    int         `json:"code"`
	Info    interface{} `json:"info"`
	ErrCode interface{} `json:"err_code"`
}

func jsonReturn(ctx *faygo.Context, code int, data interface{}, count ...interface{}) (err error) {
	var j_data interface{}
	j_data = return_jonData(code, data, count...)
	return ctx.JSON(200, j_data)
}

func return_jonData(code int, data interface{}, count ...interface{}) data_info {
	json := make(map[string]interface{})
	json["code"] = code
	json["data"] = data
	if len(count) == 1 {
		json["count"] = count[0]
	}
	return json
}

func curl_get(url string) (data jsonData, err error) {
	client := &http.Client{Timeout: time.Second * 5}
	res, err := client.Get(url)
	if err != nil {
		return
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &data)
	return
}

func curl_post(u string, param map[string]string, duration time.Duration) (data jsonData, err error) {
	client := &http.Client{
		Timeout: duration,
	}
	p := url.Values{}
	for k, v := range param {
		p[k] = []string{v}
	}
	resp, err := client.PostForm(u, p)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &data)
	return
}

//保留n位小数
func Tofix(f float64, n int) (res float64, err error) {
	format := "%." + strconv.Itoa(n) + "f"
	float_str := fmt.Sprintf(format, f)
	res, err = strconv.ParseFloat(float_str, 64)
	return
}
