package watch

import "github.com/tebeka/selenium"

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
