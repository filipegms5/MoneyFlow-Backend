package router

import (
	"github.com/filipegms5/MoneyFlow-Backend/controllers"
	"github.com/filipegms5/MoneyFlow-Backend/middlewares"
	"github.com/filipegms5/MoneyFlow-Backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	estabelecimentoController := controllers.NewEstabelecimentoController(db)
	formaPagamentoController := controllers.NewFormaPagamentoController(db)
	transacaoController := controllers.NewTransacaoController(db)
	dadosCompraController := controllers.NewDadosCompraController(db)
	usuarioController := controllers.NewUsuarioController(db)
	services.InitUsuarioService(db)

	//Rotas Publicas
	router.POST("/login", usuarioController.Login)
	router.POST("/signup", usuarioController.Create)

	//Rotas Protegidas
	auth := middlewares.AuthMiddleware()
	protected := router.Group("")

	protected.Use(auth)
	estabelecimentoRoutes := protected.Group("/estabelecimentos")
	{
		estabelecimentoRoutes.POST("", estabelecimentoController.Create)
		estabelecimentoRoutes.GET("", estabelecimentoController.GetAll)
		estabelecimentoRoutes.GET("/:id", estabelecimentoController.GetByID)
		estabelecimentoRoutes.PUT("/:id", estabelecimentoController.Update)
		estabelecimentoRoutes.DELETE("/:id", estabelecimentoController.Delete)
	}

	formaPagamentoRoutes := protected.Group("/formas-pagamento")
	{
		formaPagamentoRoutes.POST("", formaPagamentoController.Create)
		formaPagamentoRoutes.GET("", formaPagamentoController.GetAll)
		formaPagamentoRoutes.GET("/:id", formaPagamentoController.GetByID)
		formaPagamentoRoutes.PUT("/:id", formaPagamentoController.Update)
		formaPagamentoRoutes.DELETE("/:id", formaPagamentoController.Delete)
	}

	transacaoRoutes := protected.Group("/transacoes")
	{
		transacaoRoutes.POST("", transacaoController.Create)
		transacaoRoutes.GET("", transacaoController.GetAll)
		transacaoRoutes.GET("/:id", transacaoController.GetByID)
		transacaoRoutes.GET("/tipo/:tipo", transacaoController.GetByTipo)
		transacaoRoutes.PUT("/:id", transacaoController.Update)
		transacaoRoutes.DELETE("/:id", transacaoController.Delete)
	}

	usuarioRoutes := protected.Group("/usuarios")
	{
		usuarioRoutes.PUT("/:id", usuarioController.Update)
		usuarioRoutes.DELETE("/:id", usuarioController.Delete)
	}

	scanRoutes := protected.Group("/scan")
	{
		scanRoutes.POST("", dadosCompraController.FetchDadosCompra)
	}

	return router
}
