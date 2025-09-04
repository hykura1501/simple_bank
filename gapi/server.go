package gapi

import (
	"fmt"

	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/pb"
	"github.com/hykura1501/simple_bank/token"
	"github.com/hykura1501/simple_bank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	store      db.Store
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
