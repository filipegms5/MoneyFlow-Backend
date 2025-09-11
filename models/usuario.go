package models

type Usuario struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Nome    string `gorm:"type:VARCHAR(100);not null" json:"nome" binding:"required"`
	Email   string `gorm:"type:VARCHAR(100);unique;not null" json:"email" binding:"required,email"`
	Senha   string `gorm:"type:VARCHAR(255);not null" json:"senha" binding:"required"`
	IsAdmin bool   `gorm:"default:false" json:"is_admin"`
}
