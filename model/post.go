package model

import (
	"time"
)

// Post GORMで使用する投稿のModel
type Post struct {
	ID          int `gorm:"column:id;primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	UserID      string     `json:"user_id" gorm:"column:user_id"`
	Comment     string     `json:"comment" gorm:"column:comment"`
	ShopName    string     `json:"shop_name" gorm:"column:shop_name"`
	ShopAddress string     `json:"shop_address" gorm:"column:shop_address"`
	Image       string     `json:"image" gorm:"column:image"`
	Score       float64    `json:"score" gorm:"column:score"`
	// Score       string     `json:"score" gorm:"column:score"`
}
