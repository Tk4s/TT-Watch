package model

import "time"

type base struct {
	ID        int       `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"not null;column:created_at" json:"created_at"`
}
