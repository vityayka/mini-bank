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

	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, r *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}
	if violations := validateUpdateUserRequest(r); violations != nil {
		return nil, validationError(violations)
	}

	arg := db.UpdateUserParams{
		ID:       authPayload.UserID,
		FullName: pgtype.Text{String: r.GetFullName(), Valid: r.FullName != nil},
		Email:    pgtype.Text{String: r.GetEmail(), Valid: r.Email != nil},
	}

	if r.Password != nil {
		hashedPassword, err := utils.HashedPassword(r.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "couldn't hash password")
		}
		arg.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to update user: %s ", err.Error())
	}

	return &pb.UpdateUserResponse{
		User: convertUser(user),
	}, nil
}

func validateUpdateUserRequest(r *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if r.Email != nil {
		if valErr := validation.ValidateEmail(r.GetEmail()); valErr != nil {
			violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
		}
	}
	if r.FullName != nil {
		if valErr := validation.ValidateFullName(r.GetFullName()); valErr != nil {
			violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
		}
	}
	if r.Password != nil {
		if valErr := validation.ValidatePassword(r.GetPassword()); valErr != nil {
			violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
		}
	}

	return violations
}
