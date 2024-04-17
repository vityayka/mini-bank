package gapi

import (
	db "bank/db/sqlc"
	"bank/pb"
	"context"
	"database/sql"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, r *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	verifyEmail, err := server.store.GetVerifyEmail(ctx, r.GetId())

	if err != nil {
		log.Err(err).Msg("validate_email_failed")
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if verifyEmail.IsUsed {
		return nil, status.Errorf(codes.PermissionDenied, "verification code is already used")
	}

	if verifyEmail.Code != r.GetCode() {
		return nil, status.Errorf(codes.NotFound, "verification code not found")
	}

	if time.Now().After(verifyEmail.ExpiredAt) {
		return nil, status.Errorf(codes.PermissionDenied, "verification code expired")
	}

	server.store.UpdateVerifyEmails(ctx, db.UpdateVerifyEmailsParams{
		ID:     verifyEmail.ID,
		IsUsed: true,
	})

	user, err := server.store.UpdateUser(ctx, db.UpdateUserParams{
		ID: verifyEmail.UserID,
		IsVerified: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
	})

	if err != nil {
		log.Err(err).Msg("update_user_failed")
		return nil, status.Errorf(codes.Internal, "something went wrong")
	}

	return &pb.VerifyEmailResponse{
		User: convertUser(user),
	}, nil
}
