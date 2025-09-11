package models

type Transacao struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	Valor             float64          `gorm:"type:DECIMAL" json:"valor" binding:"required"`
	Data              string           `gorm:"type:DATE;not null" json:"data" binding:"required"`
	FormaPagamentoID  uint             `gorm:"column:forma_pagamento_id" json:"-"`
	FormaPagamento    *FormaPagamento  `gorm:"foreignKey:FormaPagamentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"forma_pagamento,omitempty"`
	EstabelecimentoID uint             `gorm:"column:estabelecimento_id" json:"-"`
	Estabelecimento   *Estabelecimento `gorm:"foreignKey:EstabelecimentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"estabelecimento,omitempty"`
	UsuarioID         uint             `gorm:"column:usuario_id" json:"-"`
	Usuario           *Usuario         `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"usuario,omitempty"`
}
