package api

import (
	"bank/token"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeaderName       = "authorization"
	authHeaderTypeBearer = "bearer"
	authPayloadKey       = "authPayloadKey"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(authHeaderName)
		if len(authHeader) == 0 {
			err := errors.New("authorization header is missing")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authHeaderFields := strings.Fields(authHeader)
		if len(authHeaderFields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if strings.ToLower(authHeaderFields[0]) != authHeaderTypeBearer {
			err := errors.New("unsupported auth type")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authToken := authHeaderFields[1]
		authPayload, err := tokenMaker.VerifyToken(authToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}
		ctx.Set(authPayloadKey, authPayload)
		ctx.Next()
	}
}
