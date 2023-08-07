package config

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	UserID    uint   `json:"user_id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type User struct {
	ID        uint   `gorm:"primarykey"`
	Username  string `gorm:"not null" json:"user_name"`
	Password  string `gorm:"not null" json:"password"`
	Email     string `gorm:"unique" json:"email"`
	Books     []Book `gorm:"foreignKey:UserID"` // Relação "one-to-many" com books (UserID é a chave estrangeira)
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
