package main

import (
	"TT-Watch/cmd"
	"TT-Watch/connection"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"

	"github.com/sirupsen/logrus"
)

func initConfig() {
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	viper.SetConfigName("config/app")

	if e := viper.ReadInConfig(); e != nil {
		panic(e)
	}

	env := viper.GetString("env")

	connection.InitDatabase(env)
	connection.InitRedis(env)
}

func main() {
	//	text := `[[["我只想在你心中发火","I just want to start a flame in your heart",null,null,3,null,null,[[]
	//]
	//,[[["b3ad15e7a0073e77814019b341d18493","en_zh_2019q3.md"]
	//]
	//]
	//]
	//]
	//,null,"en",null,null,null,1.0,[]
	//,[["en"]
	//,null,[1.0]
	//,["en"]
	//]
	//]
	//
	//`
	//	var data []interface{}
	//	err := json.Unmarshal([]byte(text), &data)
	//	if err != nil {
	//		fmt.Println(err)
	//	} else {
	//		datas := data[0].([]interface{})
	//		fmt.Println(len(datas), datas)
	//	}
	//	os.Exit(0)

	//reg := regexp.MustCompile(`(?i)DEFI\+\+`)
	//str := " @ERCOT_ISO  is not earning that DEFI++"
	//fmt.Println(reg.MatchString(str))
	//fmt.Println(reg.FindString(str))
	//os.Exit(0)

	initConfig()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.99999",
	})

	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	var rootCmd = &cobra.Command{}
	rootCmd.AddCommand(cmd.Watch, cmd.Coin)
	rootCmd.Execute()
}
