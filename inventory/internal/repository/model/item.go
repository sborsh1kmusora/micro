package model

import "time"

type Item struct {
	UUID string    `bson:"_id"`
	Info *ItemInfo `bson:",inline"`
}

type ItemInfo struct {
	Name      string     `bson:"name"`
	Desc      string     `bson:"desc"`
	Price     float64    `bson:"price"`
	Category  string     `bson:"category"`
	CreatedAt time.Time  `bson:"createdAt"`
	UpdatedAt *time.Time `bson:"updatedAt,omitempty"`
}
