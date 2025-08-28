package models

type Estabelecimento struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	RazaoSocial string `gorm:"type:VARCHAR(100); NOT NULL" json:"nome"`
	CNPJ        string `gorm:"type:VARCHAR(14); NOT NULL; unique" json:"cnpj"`
	Endereco    string `gorm:"type:VARCHAR(255)" json:"descricao"`
}
