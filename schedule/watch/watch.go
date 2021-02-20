package watch

import (
	"TT-Watch/library/notify/wx"
	"TT-Watch/model"
	"TT-Watch/service"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/shopspring/decimal"
	"github.com/valyala/fasthttp"

	"github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
)

type coinM struct {
	Reg *regexp.Regexp
	c   string
}

type watch struct {
	port            int
	seleniumService *selenium.Service
	agent           string
	loc             *time.Location

	influences map[string]decimal.Decimal
	coinRegexp []coinM
	lock       sync.Mutex

	wg *sync.WaitGroup

	requester      *fasthttp.Client
	gcroutineLimit chan bool
}

func NewWatch(influences []model.TwitterInfluence, coins []model.Coin) *watch {
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

	w.influences = make(map[string]decimal.Decimal)
	for _, influence := range influences {
		w.influences[influence.Influence] = decimal.Zero
	}

	w.coinRegexp = append(w.coinRegexp, coinM{
		Reg: regexp.MustCompile("(?i)coin"),
		c:   "coin",
	})

	for _, coin := range coins {
		symbol := strings.ReplaceAll(coin.Coin, "+", `\+`)
		symbol = strings.ReplaceAll(symbol, ".", `\.`)
		w.coinRegexp = append(w.coinRegexp, coinM{
			Reg: regexp.MustCompile("(?i)@" + symbol),
			c:   coin.Coin,
		})
	}

	w.wg = &sync.WaitGroup{}
	w.requester = &fasthttp.Client{}
	w.gcroutineLimit = make(chan bool, 5)

	return w
}

func (w *watch) do(influence string) {
	defer w.wg.Done()
	defer func() {
		<-w.gcroutineLimit
		logrus.Infof("End %s", influence)
	}()
	logrus.Infof("Begin %s", influence)

	db := service.GetDefaultDb()

	var articles []string
	var lastPost model.TwitterPoster

	earliestTime := decimal.Zero

	if err := db.Where(&model.TwitterPoster{}).Where("poster = ?", influence).
		Order("published_time desc").
		First(&lastPost).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			wx.SendEnterpriseWx(fmt.Sprintf("查找 %s 最新发布失败: %v", influence, err), "text")
			return
		}
	} else {
		earliestTime = decimal.NewFromInt(lastPost.PublishedTime.Unix())
		w.lock.Lock()
		w.influences[influence] = decimal.NewFromInt(lastPost.PublishedTime.Unix())
		w.lock.Unlock()
	}

	caps, _ := w.getCapabilities("chrome", "", w.agent)
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", w.port))
	if err != nil {
		logrus.Errorf("Failed to new remote, %v", err)
		wx.SendEnterpriseWx(fmt.Sprintf("Failed to new remote, %v", err), "text")
		return
	}
	defer wd.Quit()
	defer wd.DeleteAllCookies()

	_ = wd.SetPageLoadTimeout(40 * time.Second)

	if err := wd.Get("https://twitter.com/" + influence); err != nil {
		logrus.Errorf("Failed to get twitter, %v", err)
		return
	}

	saveList := []model.TwitterPoster{}

	if err := wd.WaitWithTimeout(func(wd selenium.WebDriver) (b bool, e error) {
		time.Sleep(2 * time.Second)
		wd.KeyDown(selenium.SpaceKey)
		time.Sleep(2 * time.Second)
		//wd.KeyDown(selenium.SpaceKey)
		//time.Sleep(2 * time.Second)

		//var section = $("section > div > div")  获取区域对象
		// for (var i =0;i<section.childNodes.length;i ++)  遍历子节点获取内容
		//    if (i==0) one = section.childNodes[0]
		// one.getElementsByTagName("a") 获取a标签
		// oneatitle = a[a.length-3] 获取个数-3 然后判断href 是否等于目标用户 /elonmusk

		if sections, err := wd.FindElements(selenium.ByCSSSelector, "section > div > div > div > div > div > article"); err == nil {
			// 如果没有获取到数据就重新
			if len(sections) == 0 {
				return false, nil
			}

			//if arts, err := sections.FindElements(selenium.ByCSSSelector, "div"); err == nil {
			for _, section := range sections {
				if tweets, err := section.FindElements(selenium.ByCSSSelector, "div[data-testid='tweet']"); err == nil {
					for _, tweet := range tweets {
						if aTags, err := tweet.FindElements(selenium.ByCSSSelector, "a"); err == nil {
							//if len(aTags) == 4 || len(aTags) == 5 {

							href, _ := aTags[1].GetAttribute("href")
							if strings.HasSuffix(href, influence) {
								//logrus.Infof("a tags: %d", len(aTags))
								//logrus.Infof("href: %+v", href)

								if arts, err := tweet.FindElements(selenium.ByCSSSelector, "div[dir='auto']"); err == nil {
									var content, dateStr string
									var body []string
									var datetime time.Time
									var link string

									link, _ = aTags[2].GetAttribute("href")

									if date, err := tweet.FindElement(selenium.ByCSSSelector, "time"); err == nil {
										dateStr, _ = date.GetAttribute("datetime")
										datetime, _ = time.ParseInLocation(time.RFC3339, dateStr, w.loc)
									}

									datetimeDecimal := decimal.NewFromInt(datetime.Unix())

									// 这里不用锁, 因为不会出现同时操作一个key
									if earliestTime.LessThan(datetimeDecimal) {
										earliestTime = datetimeDecimal
									}

									w.lock.Lock()
									influenceLast := w.influences[influence]
									w.lock.Unlock()

									if !influenceLast.IsZero() {
										if influenceLast.GreaterThanOrEqual(datetimeDecimal) {
											continue
										}
									}

									//if datetime.Sub(w.influences[influence]) > 0 {
									//	w.influences[influence] = datetime
									//} else {
									//	continue
									//}

									if len(arts) > 3 {
										arts = arts[3:]
										for _, art := range arts {
											text, _ := art.Text() // 只拿最后一个做内容
											text = strings.ReplaceAll(text, "\n", " ")
											text = strings.ReplaceAll(text, "\r", " ")
											text = strings.TrimLeft(text, ".")
											if text != "" {
												body = append(body, text)
											}
										}
									}

									if len(body) > 0 {
										resource := strings.Join(body, " ")
										//translate, err := w.Translate(resource)
										//if err == nil {
										//	content = fmt.Sprintf("datetime: %v, resource: %s, translate: %s", datetime.Format("2006-01-02 15:04:05"), translate["source"], translate["target"])
										//	articles = append(articles, content)
										//} else {

										contain := false

										for _, exp := range w.coinRegexp {
											if exp.Reg.MatchString(resource) {
												c := exp.Reg.FindString(resource)
												resource = strings.ReplaceAll(resource, c, "**"+c+"**")
												contain = true
												break
											}
										}

										if contain {
											resource = fmt.Sprintf("%s    ( <font color=\"info\">link</font>: [twitter detail](%s) )", resource, link)
											content = fmt.Sprintf("- <font color=\"info\">datetime</font>: %v, <font color=\"info\">post</font>: %s", datetime.Format("2006-01-02 15:04:05"), resource)
											saveList = append(saveList, model.TwitterPoster{
												Poster:        influence,
												Content:       resource,
												PublishedTime: datetime,
											})
											articles = append(articles, content)
										}

										//logrus.Warnf("Failed to translate: %v", err)
										//}

									}
								} else {
									logrus.Error(err)
								}
							}
						}
					}
				}
			}
		} else {
			return false, err
		}

		return true, nil
	}, time.Second*10); err != nil {
		logrus.Warnf("wait timeout:%s,  %v", influence, err)
	}

	if earliestTime.IsZero() {
		earliestTime = decimal.NewFromInt(time.Now().Unix())
	}

	w.lock.Lock()
	w.influences[influence] = earliestTime
	w.lock.Unlock()

	if len(articles) > 0 {
		wx.SendEnterpriseWx(fmt.Sprintf("## %s: \n\n%s", influence, strings.Join(articles, "\n\n")), "markdown")

		tx := db.Begin()
		for _, s := range saveList {
			if err := tx.Create(&s).Error; err != nil {
				tx.Rollback()
				wx.SendEnterpriseWx(fmt.Sprintf("存入数据库失败:%s %v", influence, err), "text")
				return
			}
		}

		tx.Commit()
		//logrus.Infof("%s post: \n %s", influence, strings.Join(articles, " \n"))
	}

}
