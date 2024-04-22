package api

import (
	"bank/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (server *Server) renewToken(ctx *gin.Context) {
	var request renewTokenRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := server.tokenMaker.VerifyToken(request.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if payload.ExpiresAt.Before(time.Now()) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("token expired")))
		return
	}

	session, err := server.store.GetSession(ctx, payload.ID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("blocked session")))
		return
	}

	user, err := server.store.GetUser(ctx, payload.UserID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	if session.UserID != user.ID {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("bad session")))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, utils.Role(user.Role), server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := renewTokenResponse{
		accessToken,
		accessPayload.ExpiresAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}
