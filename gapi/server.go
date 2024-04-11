package gapi

import (
	"bank/async"
	db "bank/db/sqlc"
	"bank/pb"
	"bank/token"
	"bank/utils"
)

type Server struct {
	pb.UnimplementedBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          *utils.Config
	taskDistributor async.TaskDistrubutor
}

func NewServer(config utils.Config, store db.Store, taskDistributor async.TaskDistrubutor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          &config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
