package service

import (
	"context"

	"github.com/Artem09076/lab_microservice.git/internal/proxyproto"
	sqlc "github.com/Artem09076/lab_microservice.git/internal/userdb"
	"github.com/Artem09076/lab_microservice.git/services/connect-service/internal/config"
	"github.com/Artem09076/lab_microservice.git/services/connect-service/internal/kc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	proxyproto.UnimplementedCentrifugoProxyServer
	conn     *pgxpool.Pool
	queries  *sqlc.Queries
	kcClient *kc.KCClient
}

func New(cfg *config.Config) (*Service, error) {
	connCfg, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(context.Background(), connCfg)
	if err != nil {
		return nil, err
	}
	kcClient := kc.New(cfg.KeyCloakURL, cfg.KeyCloakRealm, cfg.KeyCloakClient, cfg.KeyCloakSecret)
	return &Service{
		conn:     conn,
		queries:  sqlc.New(conn),
		kcClient: kcClient,
	}, nil
}
