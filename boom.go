package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func main() {
	log.Println("开始")

	// 启动Chrome浏览器
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless",
			"--disable-gpu",
			"--no-sandbox",
			"--disable-dev-shm-usage",
			"--blink-settings=imagesEnabled=false",
		},
	}
	caps.AddChrome(chromeCaps)
	wd, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Println("启动浏览器失败:", err)
		return
	}
	defer wd.Quit()

	// 打开网站
	err = wd.Get("http://www.txxzc.com/#/login/register")
	if err != nil {
		log.Println("打开网站失败:", err)
		return
	}

	// 循环运行100次
	for i := 0; i < 100000; i++ {
		// 输入随机手机号
		phoneNumber := generateRandomPhoneNumber()
		err = wd.WaitVisible(selenium.ByCSSSelector, `input[placeholder="请填写11位手机号"]`, 5*time.Second)
		if err != nil {
			log.Println("等待手机号输入框可见失败:", err)
			return
		}
		elem, err := wd.FindElement(selenium.ByCSSSelector, `input[placeholder="请填写11位手机号"]`)
		if err != nil {
			log.Println("查找手机号输入框失败:", err)
			return
		}
		err = elem.SendKeys(phoneNumber)
		if err != nil {
			log.Println("输入手机号失败:", err)
			return
		}
		log.Println("成功输入手机号:", phoneNumber)

		// 点击获取验证码按钮
		err = wd.WaitVisible(selenium.ByCSSSelector, "span.verification-code", 5*time.Second)
		if err != nil {
			log.Println("等待获取验证码按钮可见失败:", err)
			return
		}
		elem, err = wd.FindElement(selenium.ByCSSSelector, "span.verification-code")
		if err != nil {
			log.Println("查找获取验证码按钮失败:", err)
			return
		}
		err = elem.Click()
		if err != nil {
			log.Println("点击获取验证码按钮失败:", err)
			return
		}
		log.Println("成功点击获取验证码按钮")

		// 等待随机秒数
		ind := rand.Intn(2)
		waitDuration := time.Duration(ind) * time.Second
		log.Println("等待", ind, "秒后继续下一轮操作", i, "次")
		time.Sleep(waitDuration)
	}
}

// 生成随机手机号
func generateRandomPhoneNumber() string {
	number := rand.Intn(89999999) + 10000000 // 生成8位随机数
	return "138" + strconv.Itoa(number)
}
