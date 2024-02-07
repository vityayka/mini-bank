package api

import (
	db "bank/db/sqlc"
	"bank/utils"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type createUserResponse struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var request createUserRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashedPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashedPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := createUserResponse{
		ID:                user.ID,
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusCreated, rsp)
}
