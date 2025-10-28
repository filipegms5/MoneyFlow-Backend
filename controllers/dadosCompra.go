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

	// Bind the JSON payload to the struct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the URL variable
	url := payload.URL

	dadosCompra, err := services.FetchTransacao(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save or update Estabelecimento
	if dadosCompra.Estabelecimento != nil {
		estabelecimento, err := c.estabelecimentoRepo.GetByCnpj(dadosCompra.Estabelecimento.CNPJ)

		if estabelecimento == nil {
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

	// Save or update FormaPagamento
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

	// Attach authenticated user to the transacao
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
