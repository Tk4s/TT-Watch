package service

import (
	"database/sql/driver"

	"github.com/Tk4s/godbutils/database/sql"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func GetDefaultDb() *gorm.DB {
	var db *gorm.DB
	var er error
	for i := 0; i < 5; i++ {

		db, er = sql.GetInstanceWithName("default")
		if er != nil {
			logrus.Errorf("Failed to get db: default, %v", er)
			if er != driver.ErrBadConn {
				break
			}
		} else {
			break
		}
	}

	return db
}
