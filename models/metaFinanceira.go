package models

type MetaFinanceira struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	Valor      float64  `gorm:"type:DECIMAL" json:"valor" binding:"required"`
	DataInicio string   `gorm:"type:DATE;not null" json:"data_inicio" binding:"required"`
	DataFim    string   `gorm:"type:DATE;not null" json:"data_fim" binding:"required"`
	Descricao  string   `gorm:"type:VARCHAR(255)" json:"descricao,omitempty"`
	UsuarioID  uint     `gorm:"column:usuario_id" json:"-"`
	Usuario    *Usuario `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"usuario,omitempty"`
}
