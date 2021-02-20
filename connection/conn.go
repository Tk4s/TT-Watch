package connection

import (
	"fmt"

	"github.com/Tk4s/godbutils/database/mredis"
	"github.com/Tk4s/godbutils/database/sql"
	"github.com/spf13/viper"
)

func InitDatabase(env string) {
	sqlCfgs := viper.GetStringMap(fmt.Sprintf("database.%s", env))
	for name, value := range sqlCfgs {
		c := value.(map[string]interface{})
		cfg := sql.Config{
			Host:      c["host"].(string),
			Port:      int(c["port"].(int64)),
			UserName:  c["username"].(string),
			Password:  c["password"].(string),
			Database:  c["database"].(string),
			MaxIde:    int(c["max_ide"].(int64)),
			MaxOpen:   int(c["max_open"].(int64)),
			Charset:   c["charset"].(string),
			ParseTime: c["parse_time"].(string),
			Loc:       c["loc"].(string),
		}

		sql.InitInstanceWithName(name, cfg)
	}
}

func InitRedis(env string) {
	redisCfg := viper.Sub(fmt.Sprintf("redis.%s", env))

	mredis.NewRedis(redisCfg.GetString("address"), redisCfg.GetString("password"), redisCfg.GetInt("db"), redisCfg.GetString("network"))
}
