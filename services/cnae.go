package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/filipegms5/MoneyFlow-Backend/models"
	"gorm.io/gorm"
)

// Estrutura mínima para desserializar a resposta da BrasilAPI de CNPJ
type brasilAPICNPJResponse struct {
	CnaeFiscal interface{} `json:"cnae_fiscal"`
}

// NormalizeCNPJ remove todos os caracteres não numéricos do CNPJ
func NormalizeCNPJ(cnpj string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(cnpj, "")
}

// FetchCNAEFiscalByCNPJ consulta a BrasilAPI e retorna o CNAE fiscal principal do CNPJ
// Observação: a API pode retornar o campo como número ou string; lidamos com ambos.
func FetchCNAEFiscalByCNPJ(ctx context.Context, cnpj string) (int, error) {
	norm := NormalizeCNPJ(cnpj)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://brasilapi.com.br/api/cnpj/v1/"+norm, nil)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data brasilAPICNPJResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	// Tratamento quando o campo vem como número ou como string
	switch v := data.CnaeFiscal.(type) {
	case float64:
		return int(v), nil
	case string:
		v = strings.TrimSpace(v)
		// Remove caracteres não numéricos por segurança
		digits := NormalizeCNPJ(v)
		if len(digits) == 7 { // CNAEs principais costumam ter 7 dígitos
			// Converte para inteiro (os dígitos já são numéricos)
			var n int
			for _, ch := range digits {
				n = n*10 + int(ch-'0')
			}
			return n, nil
		}
	}
	return 0, nil
}

// EnsureCategoria garante que exista uma Categoria correspondente ao CNAE informado.
// - Se o nome mapeado for "Outros", usa/gera a categoria padrão (ID 9999999).
// - Caso contrário, usa o próprio código CNAE como ID da categoria (find-or-create).
func EnsureCategoria(db *gorm.DB, cnae int, nome string) (*models.Categoria, error) {
	// Se não houver mapeamento conhecido, usar categoria 'Outros' (ID fixo)
	if nome == "Outros" {
		outrosID := uint(9999999)
		cat := models.Categoria{ID: outrosID}
		if err := db.First(&cat, cat.ID).Error; err == nil {
			return &cat, nil
		}
		cat.Nome = "Outros"
		cat.CnaeId = NormalizeCNPJ(fmt.Sprintf("%d", outrosID))
		if err := db.Create(&cat).Error; err != nil {
			return nil, err
		}
		return &cat, nil
	}

	// ID como o próprio código CNAE
	cat := models.Categoria{ID: uint(cnae)}
	if err := db.First(&cat, cat.ID).Error; err == nil {
		return &cat, nil
	}
	cat.Nome = nome
	cat.CnaeId = NormalizeCNPJ(fmt.Sprintf("%d", cnae))
	if err := db.Create(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}
