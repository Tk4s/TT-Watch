package model

type Coin struct {
	Coin   string `json:"coin"`
	Enable int64  `json:"enable"`
}

func (Coin) TableName() string {
	return "coin"
}
