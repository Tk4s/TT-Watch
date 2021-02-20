package wx

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/valyala/fasthttp"
)

const (
	uri = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=ea147034-2755-4801-8c65-01161c62922b"
)

var client = &fasthttp.Client{}

func SendEnterpriseWx(text, typ string) {
	body := map[string]interface{}{
		"msgtype": typ,
		typ: map[string]interface{}{
			"content":        text,
			"mentioned_list": []string{"@all"},
		},
	}

	jsonBody, _ := json.Marshal(body)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetContentType("application/json")

	req.SetBody(jsonBody)

	response := fasthttp.AcquireResponse()
	if err := client.DoTimeout(req, response, 10*time.Second); err != nil {
		logrus.Error(err)
	} else {
		logrus.Infof("%+v", string(response.Body()))
	}

}
