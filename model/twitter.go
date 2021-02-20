package model

import "time"

type TwitterPoster struct {
	base
	Poster        string    `json:"poster"`
	Content       string    `json:"content"`
	PublishedTime time.Time `json:"published_time"`
}

func (TwitterPoster) TableName() string {
	return "twitter_poster"
}

type TwitterInfluence struct {
	base
	Influence string `json:"influence"`
	Enable    int64  `json:"enable"`
}

func (TwitterInfluence) TableName() string {
	return "twitter_influence"
}
