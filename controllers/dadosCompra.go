package controllers

import (
	"net/http"

	"github.com/filipegms5/MoneyFlow-Backend/repositories"
	"github.com/filipegms5/MoneyFlow-Backend/services"
	"github.com/filipegms5/MoneyFlow-Backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DadosCompraController struct {
	transacaoRepo       *repositories.TransacaoRepository
	estabelecimentoRepo *repositories.EstabelecimentoRepository
	formaPagamentoRepo  *repositories.FormaPagamentoRepository
}

type RequestPayload struct {
	URL string `json:"url" binding:"required"`
}

func NewDadosCompraController(db *gorm.DB) *DadosCompraController {
	return &DadosCompraController{
		transacaoRepo:       repositories.NewTransacaoRepository(db),
		estabelecimentoRepo: repositories.NewEstabelecimentoRepository(db),
		formaPagamentoRepo:  repositories.NewFormaPagamentoRepository(db),
	}
}

func (c *DadosCompraController) FetchDadosCompra(ctx *gin.Context) {
	var payload RequestPayload

	// Faz o bind do JSON para a struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extrai a URL do payload
	url := payload.URL

	dadosCompra, err := services.FetchTransacao(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Salva ou atualiza o Estabelecimento
	if dadosCompra.Estabelecimento != nil {
		estabelecimento, err := c.estabelecimentoRepo.GetByCnpj(dadosCompra.Estabelecimento.CNPJ)

		if estabelecimento == nil {
			// Se não houver CNPJ, define categoria "Outros" automaticamente
			if dadosCompra.Estabelecimento.CNPJ == "" && dadosCompra.Estabelecimento.CategoriaID == 0 {
				dadosCompra.Estabelecimento.CategoriaID = 9999999
			}
			err = c.estabelecimentoRepo.Create(dadosCompra.Estabelecimento)
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if estabelecimento != nil {
			dadosCompra.Estabelecimento = estabelecimento
		}

		dadosCompra.EstabelecimentoID = dadosCompra.Estabelecimento.ID
		dadosCompra.Estabelecimento = nil
	}

	// Salva ou atualiza a Forma de Pagamento
	if dadosCompra.FormaPagamento != nil {
		formaPagamento, err := c.formaPagamentoRepo.GetByNome(dadosCompra.FormaPagamento.Nome)

		if formaPagamento == nil {
			err = c.formaPagamentoRepo.Create(dadosCompra.FormaPagamento)
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if formaPagamento != nil {
			dadosCompra.FormaPagamento = formaPagamento
		}

		dadosCompra.FormaPagamentoID = dadosCompra.FormaPagamento.ID
		dadosCompra.FormaPagamento = nil
	}

	// Anexa o usuário autenticado à transação
	if uid, ok := utils.GetUserIDFromContext(ctx); ok {
		dadosCompra.UsuarioID = uid
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in token"})
		return
	}

	err = c.transacaoRepo.Create(&dadosCompra)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dadosCompra)
}
