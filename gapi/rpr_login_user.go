package gapi

import (
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUserByEmail(ctx, r.GetEmail())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user is not found")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = utils.CompareHashAndPassword(user.HashedPassword, r.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, "wrong password")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    time.Now().Add(server.config.RefreshTokenDuration),
	})

	return &pb.LoginUserResponse{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt),
		User:                  convertUser(user),
	}, nil
}
