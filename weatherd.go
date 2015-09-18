package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"flag"
	"strconv"
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

var appid = flag.String("appid", "", "the appid")
var privkey = flag.String("privkey", "", "the privkey")
var areaids = flag.String("areaids", "", "the areaids , more ids use | to split")
var port = flag.Int("port", 9090, "listen port")
var span = flag.Int("span", 15, "interval of get weather from official site(minute)")

var curWeather string
var coder = base64.StdEncoding

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func hmacenc(data []byte, key []byte) []byte {

	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
func getWeather() {

	t := time.Now()

	for {
		pubkey := "http://open.weather.com.cn/data/?areaid=" + *areaids + "&type=forecast_v&date=" + t.Format("200601021504") + "&appid=" + *appid
		key := base64Encode(hmacenc([]byte(pubkey), []byte(*privkey)))
		rep, err := httpGet(pubkey[:len(pubkey)-(len(*appid)-6)] + "&key=" + url.QueryEscape(string(key)))
		if err == nil {

			curWeather = rep
		} else {
			print(err.Error())
		}
		<-time.After(15 * time.Minute)

	}

}

func procReq(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(curWeather)))

	fmt.Fprintf(w, curWeather)

}

func main() {
	flag.Parse()
	if *appid == "" {
		println("please provide appid")
		return
	}
	if *privkey == "" {
		println("please provide privkey")
		return
	}
	if *areaids == "" {
		println("please provide city areaids")
		return
	}
	// 设置访问的路由
	go getWeather()
	http.HandleFunc("/", procReq)

	// 设置监听的端口
	err := http.ListenAndServe(":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
