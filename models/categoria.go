package models

type Categoria struct {
	ID     uint   `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Nome   string `gorm:"type:VARCHAR(100); NOT NULL" json:"nome"`
	CnaeId string `gorm:"type:VARCHAR(10); NOT NULL" json:"cnae_id"`
}
