package model

import "time"

type Item struct {
	UUID string
	Info *ItemInfo
}

type ItemInfo struct {
	Name      string
	Desc      string
	Price     float64
	Category  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
