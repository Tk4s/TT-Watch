package watch

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

type watch struct {
	port            int
	seleniumService *selenium.Service
	agent           string
}

func NewWatch() *watch {
	w := &watch{}
	err := w.initPort()
	if err != nil {
		return nil
	}

	err = w.initService()
	if err != nil {
		return nil
	}

	w.agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.58 Safari/537.36"

	return w
}

func (w *watch) do(star string, wg sync.WaitGroup) {
	defer wg.Done()

	caps, _ := w.getCapabilities("chrome", "", w.agent)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", w.port))
	if err != nil {
		logrus.Errorf("Failed to new remote, %v", err)
		return
	}
	defer wd.Quit()
	defer wd.DeleteAllCookies()

	_ = wd.SetPageLoadTimeout(40 * time.Second)

	if err := wd.Get("https://twitter.com/elonmusk"); err != nil {
		logrus.Errorf("Failed to get twitter, %v", err)
		return
	}

	time.Sleep(2 * time.Second)
	wd.KeyDown(selenium.SpaceKey)
	time.Sleep(2 * time.Second)
	wd.KeyDown(selenium.SpaceKey)
	time.Sleep(2 * time.Second)

	//var section = $("section > div > div")  获取区域对象
	// for (var i =0;i<section.childNodes.length;i ++)  遍历子节点获取内容
	//    if (i==0) one = section.childNodes[0]
	// one.getElementsByTagName("a") 获取a标签
	// oneatitle = a[a.length-3] 获取个数-3 然后判断href 是否等于目标用户 /elonmusk

	if sections, err := wd.FindElements(selenium.ByCSSSelector, "section > div > div > div"); err == nil {
		//if arts, err := sections.FindElements(selenium.ByCSSSelector, "div"); err == nil {
		for _, section := range sections {
			if tweets, err := section.FindElements(selenium.ByCSSSelector, "div[data-testid='tweet']"); err == nil {
				for _, tweet := range tweets {
					if aTags, err := tweet.FindElements(selenium.ByCSSSelector, "a"); err == nil {
						//if len(aTags) == 4 || len(aTags) == 5 {

						href, _ := aTags[len(aTags)-3].GetAttribute("href")
						if strings.Contains(href, "/elonmusk") {
							logrus.Infof("a tags: %d", len(aTags))
							logrus.Infof("href: %+v", href)

							if arts, err := tweet.FindElements(selenium.ByCSSSelector, "div[dir='auto']"); err == nil {
								if len(arts) > 3 {
									text, err := arts[len(arts)-1].Text() // 只拿最后一个做内容
									text = strings.ReplaceAll(text, "\n", " ")
									text = strings.ReplaceAll(text, "\r", " ")
									logrus.Infof("contnet: %+v, %v", text, err)
								}

							} else {
								logrus.Error(err)
							}
							fmt.Println("")
						}
					}
				}
			}
		}
	}
}