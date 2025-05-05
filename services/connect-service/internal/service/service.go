package service

import (
	"context"
	"s0709-22/internal/proxyproto"
	"s0709-22/internal/userdb"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Service ...
type Service struct {
	proxyproto.UnimplementedCentrifugoProxyServer
	conn    *pgxpool.Pool
	storage *userdb.Queries
}

// New ...
func New(uri string) (*Service, error) {
	connCfg, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), connCfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		conn:    conn,
		storage: userdb.New(conn),
	}, nil
}
