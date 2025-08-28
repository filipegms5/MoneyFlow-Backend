package models

type FormaPagamento struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Nome string `gorm:"type:VARCHAR(100); NOT NULL" json:"nome"`
}
