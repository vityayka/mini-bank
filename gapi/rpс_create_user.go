package gapi

import (
	"bank/async"
	db "bank/db/sqlc"
	"bank/pb"
	"bank/utils"
	"bank/validation"
	"context"
	"errors"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, r *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(r); violations != nil {
		return nil, validationError(violations)
	}
	hashedPassword, err := utils.HashedPassword(r.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "couldn't hash password")
	}

	arg := db.CreateUserParams{
		Username:       r.GetUsername(),
		Role:           string(utils.Depositor),
		HashedPassword: hashedPassword,
		FullName:       r.GetFullName(),
		Email:          r.GetEmail(),
	}

	result, err := server.store.CreateUserTX(ctx, db.CreateUserTxParams{
		CreateUserParams: arg,
		AfterCreate: func(user db.User) error {
			payload := &async.PayloadSendVerifyEmail{UserID: user.ID}
			opts := []asynq.Option{
				asynq.ProcessIn(10 * time.Second),
				asynq.MaxRetry(5),
				// asynq.Queue("critical"),
			}
			if err = server.taskDistributor.DistributeTaskVerifyEmail(ctx, payload, opts...); err != nil {
				return status.Errorf(codes.Internal, err.Error())
			}

			return nil
		},
	})

	if err != nil {
		if errors.Is(err, db.ErrUserAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		log.Println(err)
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		User: convertUser(result.User),
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
