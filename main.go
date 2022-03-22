package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var cfgMap = make(map[string]string)
var fileName = "app.txt"
var regexPattern = `(.*?)=(.*)`

func FilterConfigAndCallFn(prefix string, config map[string]string, fn func(k, v string)) {
	for k, v := range config {
		if strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
			fn(k, v)
		}
	}
}
func ReadConfig() {
	fileBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r, _ := regexp.Compile(regexPattern)
	for _, match := range r.FindAllStringSubmatch(string(fileBytes), -1) {
		//fmt.Println(match)
		if !strings.HasPrefix(match[1], "#") {
			cfgMap[match[1]] = match[2]
		}
	}
}
func printConfig() {
	fmt.Println("Config details")
	for k, v := range cfgMap {
		fmt.Printf("%s : \"%s\"\n", k, v)
	}
}

func main() {
	// Read Config File app.env
	ReadConfig()
	// Print Config details
	printConfig()

	// check ignore bad ssl setting
	if val, ok := cfgMap["app_ignore_bad_ssl_error"]; ok && val == "1" {
		fmt.Println("Ignoring bad ssl error: ON")
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if checkLogin() {
		fmt.Println("You're logged on to BT Wi-Fi")
	} else {
		login()
	}

	if val, ok := cfgMap["app_login_check_timer"]; ok && val != "" {
		checkTimerSeconds, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			fmt.Println("Invalid app_login_check_timer value. It should be a number")
			fmt.Println(err)
			return
		}
		fmt.Printf("Setting up background login check every %d seconds\n", checkTimerSeconds)

		checkLoginTimer(time.Duration(checkTimerSeconds) * time.Second)
	}

}
func checkLoginTimer(duration time.Duration) {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			if checkLogin() {
				fmt.Println("You're logged on to BT Wi-Fi")
				continue
			}
			login()
		case <-interrupt:
			break loop
		}
	}
}
func login() bool {
	params := url.Values{}
	FilterConfigAndCallFn("post_", cfgMap, params.Add)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", cfgMap["app_login_url"], body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")

	FilterConfigAndCallFn("header_", cfgMap, req.Header.Set)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	if checkLogin() {
		fmt.Println("You're logged on to BT Wi-Fi")
		return true
	}

	fmt.Println("There is someting went wrong")
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return false
	}
	bodyString := string(bodyBytes)
	fmt.Println(bodyString)
	return false
}
func checkLogin() bool {

	req, err := http.NewRequest("GET", cfgMap["app_login_check_url"], nil)
	if err != nil {
		fmt.Println(err)
		return false
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.93 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	FilterConfigAndCallFn("header_", cfgMap, req.Header.Set)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		bodyString := string(bodyBytes)
		// fmt.Println(bodyString)
		if strings.Contains(bodyString, cfgMap["app_login_check_keyword"]) {
			return true
		} else {
			fmt.Println("You are not logged on to BT")
		}
	} else {
		fmt.Println("Server Response Status Code", resp.StatusCode)
	}

	return false
}
