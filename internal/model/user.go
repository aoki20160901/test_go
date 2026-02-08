package model

type User struct {
	ID   string `json:"id" gorm:"primaryKey;type:text"`
	Name string `json:"name" gorm:"not null"`
}

// テーブル名を明示（任意）
func (User) TableName() string {
	return "users"
}
