package models

type Transacao struct {
	ID                uint             `gorm:"primaryKey" json:"id"`
	Valor             float64          `gorm:"type:DECIMAL" json:"valor" binding:"required"`
	Data              string           `gorm:"type:DATE;not null" json:"data" binding:"required"`
	Tipo              string           `gorm:"type:VARCHAR(50);not null;default:despesa" json:"tipo" binding:"required,oneof=despesa receita"`
	Descricao         string           `gorm:"type:VARCHAR(255)" json:"descricao,omitempty"`
	Recorrente        bool             `gorm:"type:BOOLEAN;not null;default:false" json:"recorrente,omitempty"`
	FormaPagamentoID  uint             `gorm:"column:forma_pagamento_id" json:"forma_pagamento_id,omitempty"`
	FormaPagamento    *FormaPagamento  `gorm:"foreignKey:FormaPagamentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"forma_pagamento,omitempty"`
	EstabelecimentoID uint             `gorm:"column:estabelecimento_id" json:"estabelecimento_id,omitempty"`
	Estabelecimento   *Estabelecimento `gorm:"foreignKey:EstabelecimentoID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"estabelecimento,omitempty"`
	UsuarioID         uint             `gorm:"column:usuario_id" json:"-"`
	Usuario           *Usuario         `gorm:"foreignKey:UsuarioID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"usuario,omitempty"`
	CategoriaID       *uint            `gorm:"column:categoria_id" json:"categoria_id,omitempty"`
	Categoria         *Categoria       `json:"categoria,omitempty"`
}
