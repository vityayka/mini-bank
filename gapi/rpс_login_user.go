package gapi

import (
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"bank/validation"
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, r *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserRequest(r); violations != nil {
		return nil, validationError(violations)
	}
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

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.ID, utils.Role(user.Role), server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.ID, utils.Role(user.Role), server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	meta := server.extractMedadata(ctx)
	server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    meta.UserAgent,
		ClientIp:     meta.ClientIP,
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

func validateLoginUserRequest(r *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if valErr := validation.ValidateEmail(r.GetEmail()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}
	if valErr := validation.ValidatePassword(r.GetPassword()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}
	return violations
}
