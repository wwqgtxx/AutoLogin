package main

import (
	"log"
	"os"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
	"net/url"
	"net/http/cookiejar"
	"time"
)

var user_name = ""
var pass_word = ""

var Logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
var accept = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8"
var accept_Encoding = "gzip, deflate"
var accept_Language = "zh-CN,zh;q=0.8"
var userAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/53.0.2785.104 Safari/537.36 Core/1.53.2669.400 QQBrowser/9.6.10990.400"

func do_login(client *http.Client) {
	Logger.Println("Detecting...")
	req, err := http.NewRequest("GET", "http://10.10.1.96/eportal/InterFace.do?method=getOnlineUserInfo", strings.NewReader(""))
	if err != nil {
		Logger.Println(err)
		return
	}
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Encoding", accept_Encoding)
	req.Header.Set("Accept-Language", accept_Language)
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		Logger.Println(err)
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logger.Println(err)
		return
	}
	var json_data map[string]interface{}
	err = json.Unmarshal(b, &json_data)
	if err != nil {
		Logger.Println(err)
	}
	if value, ok := json_data["result"]; ok {
		if result, ok := value.(string); ok {
			if result == "success" {
				Logger.Println("Was login!")
				return
			}
			Logger.Println("Detect logout!")
			Logger.Println("Try to login!")
			req, err := http.NewRequest("GET", "http://123.123.123.123", strings.NewReader(""))
			if err != nil {
				Logger.Println(err)
				return
			}
			req.Header.Set("Accept", accept)
			req.Header.Set("Accept-Encoding", accept_Encoding)
			req.Header.Set("Accept-Language", accept_Language)
			req.Header.Set("User-Agent", userAgent)
			resp, err := client.Do(req)
			if err != nil {
				Logger.Println(err)
				return
			}
			defer resp.Body.Close()
			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				Logger.Println(err)
				return
			}
			str := string(b)
			if strings.Contains(str, "http://10.10.1.96/eportal/index.jsp") {
				eportal_url := strings.Replace(strings.Replace(string(b), "<script>self.location.href='", "", -1), "'</script>", "", -1)
				queryString := url.QueryEscape(strings.Replace(eportal_url, "http://10.10.1.96/eportal/index.jsp?", "", -1))
				//Logger.Println(queryString)

				req, err = http.NewRequest("GET", "http://10.10.1.96/eportal/InterFace.do?method=pageInfo", strings.NewReader(""))
				if err != nil {
					Logger.Println(err)
					return
				}
				req.Header.Set("Accept", accept)
				req.Header.Set("Accept-Encoding", accept_Encoding)
				req.Header.Set("Accept-Language", accept_Language)
				req.Header.Set("User-Agent", userAgent)
				resp, err := client.Do(req)
				if err != nil {
					Logger.Println(err)
					return
				}
				defer resp.Body.Close()
				b, err = ioutil.ReadAll(resp.Body)
				if err != nil {
					Logger.Println(err)
					return
				}
				//Logger.Println(string(b))

				v := url.Values{"userId": {user_name}, "password": {pass_word}, "service": {""}, "queryString": {queryString},
					"operatorPwd":        {""}, "operatorUserId": {""}, "validcode": {""}}
				//Logger.Println(v.Encode())
				body := strings.NewReader(v.Encode())

				req, err = http.NewRequest("POST", "http://10.10.1.96/eportal/InterFace.do?method=login", body)
				if err != nil {
					Logger.Println(err)
				}
				req.Header.Set("Accept", accept)
				req.Header.Set("Accept-Encoding", accept_Encoding)
				req.Header.Set("Accept-Language", accept_Language)
				req.Header.Set("User-Agent", userAgent)
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

				resp, err = client.Do(req)
				if err != nil {
					Logger.Println(err)
					return
				}
				defer resp.Body.Close()
				b, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					Logger.Println(err)
					return
				}
				//Logger.Println(string(b))
				var json_data map[string]interface{}
				err = json.Unmarshal(b, &json_data)
				if err != nil {
					Logger.Println(err)
					return
				}
				if value, ok := json_data["result"]; ok {
					if result, ok := value.(string); ok {
						if result != "success" {
							Logger.Println("Login failed!")
							return
						}
					}
				}

				for ; ; {
					req, err := http.NewRequest("GET", "http://10.10.1.96/eportal/InterFace.do?method=getOnlineUserInfo", strings.NewReader(""))
					if err != nil {
						Logger.Println(err)
						return
					}
					req.Header.Set("Accept", accept)
					req.Header.Set("Accept-Encoding", accept_Encoding)
					req.Header.Set("Accept-Language", accept_Language)
					req.Header.Set("User-Agent", userAgent)
					resp, err := client.Do(req)
					if err != nil {
						Logger.Println(err)
						return
					}
					defer resp.Body.Close()
					b, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						Logger.Println(err)
						return
					}
					var json_data map[string]interface{}
					err = json.Unmarshal(b, &json_data)
					if err != nil {
						Logger.Println(err)
						return
					}
					if value, ok := json_data["result"]; ok {
						if result, ok := value.(string); ok {
							if result == "success" {
								Logger.Println("Finish login!")
								return
							}
						}

					}

				}
			}

		}
	}
}

func main() {
	config := NewConfig()
	config.Load("config.json")
	user_name = config.UserName
	pass_word = config.PassWord
	Logger.Println("Init Username:"+user_name)
	Logger.Println("Init Password:"+pass_word)
	jar, err := cookiejar.New(nil)
	if err != nil {
		Logger.Println(err)
	}
	client := &http.Client{Jar: jar}
	for ; ; {
		do_login(client)
		time.Sleep(1 * time.Second)
	}
}
