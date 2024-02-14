package api

import (
	db "bank/db/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required"`
	ToAccountID   int64  `json:"to_account_id" binding:"required"`
	Currency      string `json:"currency" binding:"required,currency"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request createTransferRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	userID, _ := getUserIdFromAuth(ctx)
	_, err = server.store.GetUserAccount(ctx, db.GetUserAccountParams{userID, request.FromAccountID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusForbidden, errorResponse(errors.New("forbidden")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	fromAccount, err := server.validAccount(ctx, request.FromAccountID, request.Currency)
	if err != nil {
		return
	}

	toAccount, err := server.validAccount(ctx, request.ToAccountID, request.Currency)
	if err != nil {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        request.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (*db.Account, error) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return nil, err
		}

		return nil, err
	}

	if account.Currency != currency {
		err = fmt.Errorf("unsupported currency %s for account %s", currency, account.Currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return nil, err
	}
	return &account, nil
}
