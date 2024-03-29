package gapi

import (
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"bank/validation"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(r); violations != nil {
		badRequest := &errdetails.BadRequest{FieldViolations: violations}
		statusInvalid := status.New(codes.InvalidArgument, "invalid params")

		statusDetails, err := statusInvalid.WithDetails(badRequest)
		if err != nil {
			return nil, statusInvalid.Err()
		}
		return nil, statusDetails
	}
	hashedPassword, err := utils.HashedPassword(r.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't hash password")
	}

	arg := db.CreateUserParams{
		Username:       r.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       r.GetFullName(),
		Email:          r.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		User: convertUser(user),
	}, nil
}

func validateCreateUserRequest(r *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if valErr := validation.ValidateEmail(r.GetEmail()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}
	if valErr := validation.ValidateFullName(r.GetFullName()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}
	if valErr := validation.ValidatePassword(r.GetPassword()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}
	if valErr := validation.ValidateUsername(r.GetUsername()); valErr != nil {
		violations = append(violations, fieldViolation(valErr.Field, valErr.Error))
	}

	return violations
}
