package watch

import (
	"fmt"
	"net"
	"time"

	"github.com/tebeka/selenium/chrome"

	"github.com/tebeka/selenium"

	"github.com/sirupsen/logrus"
)

func (w *watch) getChromeService(port int, opts []selenium.ServiceOption) (*selenium.Service, error) {
	service, err := selenium.NewChromeDriverService("chromedriver", port, opts...)
	return service, err
}

func (w *watch) initService() error {
	opts := []selenium.ServiceOption{}

	selenium.SetDebug(false)
	service, err := w.getChromeService(w.port, opts)
	if err != nil {
		logrus.Errorf("Failed to get chrome service, %v", err)
		return err
	}

	w.seleniumService = service
	return nil
}

func (w *watch) initPort() error {
	port, err := w.pickUnusedPort()
	if err != nil {
		logrus.Errorf("Failed to get unused port, err: %v", err)
		return err
	}

	w.port = port
	return nil
}

func (w *watch) pickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func (w *watch) close() {
	time.Sleep(3 * time.Second)
	w.seleniumService.Stop()
}

func (w *watch) getCapabilities(webDriver string, proxy string, agent string) (selenium.Capabilities, string) {
	imgCaps := map[string]interface{}{
		//"profile.managed_default_content_settings.images": 2,
		//"profile.managed_default_content_settings.javascript": 2,
	}
	var caps selenium.Capabilities

	caps = selenium.Capabilities{
		"browserName": "chrome",
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			"--headless",
			"--start-maximized",
			"--no-sandbox",
			fmt.Sprintf("--user-agent=%s", agent),
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
			"--auto-open-devtools-for-tab",
		},
	}

	caps.AddChrome(chromeCaps)

	return caps, agent
}
