package models

import "time"

type Order struct {
	OrderId      uint   `gorm:"primaryKey"`
	CustomerName string `gorm:"not null;type:varchar(191)"`
	Items        []Item
	OrderedAt    time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
