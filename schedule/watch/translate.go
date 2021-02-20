package watch

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/valyala/fasthttp"
)

const (
	apiURI = "https://translation.googleapis.com/language/translate/v2"
)

var uri = "https://translate.googleapis.com/translate_a/single"

func (w *watch) Translate(text string) (map[string]string, error) {
	result := make(map[string]string)

	query := url.Values{}
	query.Add("client", "gtx")
	query.Add("sl", "auto")
	query.Add("tl", "zh-CN")
	query.Add("dt", "t")
	query.Add("q", text)

	req := fasthttp.AcquireRequest()
	req.SetRequestURI(uri)
	req.Header.Set("User-Agent", w.agent)
	req.Header.SetMethod(http.MethodGet)
	req.URI().SetQueryString(query.Encode())

	response := fasthttp.AcquireResponse()
	err := w.requester.DoTimeout(req, response, 10*time.Second)

	if err != nil {
		return result, err
	}

	if response.StatusCode() != http.StatusOK {
		return result, errors.Errorf("Translate response code is error: %d", response.StatusCode())
	}

	var source, target string

	var body []interface{}
	if err := json.Unmarshal(response.Body(), &body); err == nil && body != nil {
		datas := body[0].([]interface{})
		for _, data := range datas {
			ret := data.([]interface{})
			target = target + ret[0].(string)
			source = source + ret[1].(string)
		}
	} else {
		logrus.Error(err)
	}

	result = map[string]string{
		"source": source,
		"target": target,
	}

	return result, nil
}
