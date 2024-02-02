package api

import (
	db "bank/db/sqlc"
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var request createAccountRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    request.Owner,
		Balance:  0,
		Currency: request.Currency,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	err := ctx.ShouldBindUri(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type listAccountRequest struct {
	PageNum int32 `form:"page_num" binding:"required,min=1"`
	Count   int32 `form:"count" binding:"required,min=1,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var request listAccountRequest
	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  request.Count,
		Offset: (request.PageNum - 1) * request.Count,
	}

	accounts, err := server.store.ListAccounts(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
