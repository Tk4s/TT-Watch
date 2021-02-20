package coin

import (
	"TT-Watch/model"
	"TT-Watch/service"
	"encoding/json"

	"github.com/gogf/gf/os/gcron"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
)

type coinMarketCapModel struct {
	Data   coinMarketCapDataModel   `json:"data"`
	Status coinMarketCapStatusModel `json:"status"`
}

type coinMarketCapStatusModel struct {
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type coinMarketCapDataModel struct {
	CryptoCurrencyMap []coinMarketCapDataCryptoModel `json:"cryptoCurrencyMap"`
}

type coinMarketCapDataCryptoModel struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	IsActive int64  `json:"is_active"`
	Status   string `json:"status"`
}

func Run(cmd *cobra.Command, args []string) {
	gcron.Add("@every 1h", func() {
		db := service.GetDefaultDb()
		uri := "https://api.coinmarketcap.com/data-api/v3/map/all?listing_status=active,untracked"

		_, resp, err := fasthttp.Get(nil, uri)
		if err == nil {
			data := &coinMarketCapModel{}
			if err = json.Unmarshal(resp, data); err == nil {
				if data.Status.ErrorCode != "0" {
					logrus.Errorf("Get Coin error: %+v", data.Status)
					return
				}
				for _, d := range data.Data.CryptoCurrencyMap {
					exists := 0
					err = db.Model(&model.Coin{}).Where("coin = ?", d.Symbol).Count(&exists).Error
					if exists == 0 {
						coin := model.Coin{
							Coin:   d.Symbol,
							Enable: d.IsActive,
						}

						db.Create(&coin)
					} else {
						db.Model(&model.Coin{}).Where("coin = ?", d.Symbol).Updates(map[string]interface{}{
							"enable": d.IsActive,
						})
					}
				}
			} else {
				logrus.Error(err)
			}
		} else {
			logrus.Error(err)
		}
	})
	select {}
}
