package api

import (
	db "bank/db/sqlc"
	"bank/token"
	"bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     *utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     &config,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	if validator, isOk := binding.Validator.Engine().(*validator.Validate); isOk {
		validator.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew", server.renewToken)

	auth := router.Group("/", authMiddleware(server.tokenMaker))
	auth.POST("/accounts", server.createAccount)
	auth.GET("/accounts/:id", server.getAccount)
	auth.GET("/accounts", server.listAccounts)
	auth.POST("/transfers", server.createTransfer)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) Serve(addr string) error {
	return server.router.Run(addr)
}
