package service

import (
	"context"
	"time"

	"github.com/Artem09076/lab_microservice.git/internal/proxyproto"
	sqlc "github.com/Artem09076/lab_microservice.git/internal/userdb"
	"github.com/Artem09076/lab_microservice.git/services/connect-service/internal/config"
	"github.com/Nerzal/gocloak/v13"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	proxyproto.UnimplementedCentrifugoProxyServer
	conn           *pgxpool.Pool
	queries        *sqlc.Queries
	kcClient       *gocloak.GoCloak
	token          *gocloak.JWT
	KeyCloakRealm  string
	KeyCloakClient string
	KeyCloakSecret string
	expiredAt      time.Time
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
	kcClient := gocloak.NewClient(cfg.KeyCloakURL)

	return &Service{
		conn:           conn,
		queries:        sqlc.New(conn),
		kcClient:       kcClient,
		KeyCloakRealm:  cfg.KeyCloakRealm,
		KeyCloakClient: cfg.KeyCloakClient,
		KeyCloakSecret: cfg.KeyCloakSecret,
	}, nil
}
