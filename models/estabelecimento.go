package models

type Estabelecimento struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	RazaoSocial string     `gorm:"type:VARCHAR(100); NOT NULL" json:"nome"`
	CNPJ        string     `gorm:"type:VARCHAR(14); NOT NULL; unique" json:"cnpj"`
	Endereco    string     `gorm:"type:VARCHAR(255)" json:"descricao"`
	CategoriaID uint       `gorm:"column:categoria_id" json:"-"`
	Categoria   *Categoria `gorm:"foreignKey:CategoriaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"categoria,omitempty"`
}
