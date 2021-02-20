package watch

import (
	"TT-Watch/library/notify/wx"
	"TT-Watch/model"
	"TT-Watch/service"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gogf/gf/os/gtimer"
	"github.com/spf13/cobra"
)

//马斯克 https://twitter.com/elonmusk
//V神 https://twitter.com/VitalikButerin
//cz    https://twitter.com/cz_binance
//孙宇晨 https://twitter.com/justinsuntron
//ac    https://twitter.com/AndreCronjeTech
//SBF  https://twitter.com/SBF_Alameda
//tyler https://twitter.com/tyler
//Barry Silbert https://twitter.com/BarrySilbert
//Paolo Ardoino https://twitter.com/paoloardoino
//Michael Saylor https://twitter.com/michael_saylor
//Gavin Wood https://twitter.com/gavofyork
//Kris https://twitter.com/Kris_HK
//Brian Armstrong https://twitter.com/brian_armstrong
//John McAfee https://twitter.com/officialmcafee
//jack https://twitter.com/jack
//Roger Ver https://twitter.com/rogerkver
//Hayden Adams https://twitter.com/haydenzadams
//Charles Hoskinson https://twitter.com/IOHK_Charles

func Run(cmd *cobra.Command, args []string) {

	gtimer.AddSingleton(30*time.Second, func() {
		db := service.GetDefaultDb()
		influences := []model.TwitterInfluence{}

		if err := db.Model(&model.TwitterInfluence{}).Where("enable = 1").Find(&influences).Error; err != nil {
			logrus.Error(err)
			return
		}

		coins := []model.Coin{}

		if err := db.Model(&model.Coin{}).Where("enable = 1").Find(&coins).Error; err != nil {
			logrus.Error(err)
			return
		}

		if len(influences) == 0 || len(coins) == 0 {
			wx.SendEnterpriseWx(fmt.Sprintf("influences: %d, coin: %d.", len(influences), len(coins)), "text")
			return
		}

		w := NewWatch(influences, coins)
		if w == nil {
			return
		}

		defer w.close()

		w.wg.Add(len(w.influences))
		for influence := range w.influences {
			w.gcroutineLimit <- false
			go w.do(influence)
		}
		w.wg.Wait()

		logrus.Info("End all")
		time.Sleep(30 * time.Second)
	})

	select {}
}
