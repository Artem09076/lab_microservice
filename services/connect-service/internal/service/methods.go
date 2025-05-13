package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Artem09076/lab_microservice.git/internal/proxyproto"
	sqlc "github.com/Artem09076/lab_microservice.git/internal/userdb"
	"github.com/Nerzal/gocloak/v13"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) fetchKeycloakUser(c context.Context, userId uuid.UUID) (sqlc.User, error) {
	log.Println("fetchKeycloakUser")
	userKC, err := s.GetCloakByUser(c, userId.String())
	if err != nil {
		return sqlc.User{}, err
	}
	user := sqlc.User{
		ID:         pgtype.UUID{Valid: true, Bytes: userId},
		Username:   *userKC.Username,
		GivenName:  *userKC.FirstName,
		FamilyName: *userKC.LastName,
		Enabled:    *userKC.Enabled,
	}

	if err = s.queries.CreateUser(c, sqlc.CreateUserParams{
		ID:         pgtype.UUID{Valid: true, Bytes: userId},
		Username:   *userKC.Username,
		GivenName:  *userKC.FirstName,
		FamilyName: *userKC.LastName,
		Enabled:    *userKC.Enabled,
	}); err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (s *Service) Subscribe(c context.Context, req *proxyproto.SubscribeRequest) (*proxyproto.SubscribeResponse, error) {
	userID, err := uuid.Parse(req.User)
	log.Println("subsribe request")
	if err != nil {
		return SubscribeResponseError(107, "invalid id")
	}
	user, err := s.queries.GetUserByID(c, pgtype.UUID{Valid: true, Bytes: userID})
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.fetchKeycloakUser(c, userID)
		if err != nil {
			return SubscribeResponseError(100, "Internal server error")
		}
	} else if err != nil {
		return SubscribeResponseError(100, "Internal server error")
	}
	res, err := s.queries.UserCanSubscribe(c, sqlc.UserCanSubscribeParams{
		ID:      user.ID,
		Channel: req.Channel,
	})

	if err != nil {
		return SubscribeResponseError(100, "Internal server error")
	}

	if res == 0 {
		return SubscribeResponseError(103, "permission denied")
	}
	return &proxyproto.SubscribeResponse{}, nil
}

func (s *Service) Publish(c context.Context, req *proxyproto.PublishRequest) (*proxyproto.PublishResponse, error) {
	log.Println("Publish")
	userID, err := uuid.Parse(req.User)
	if err != nil {
		return PublishResponseError(107, "invalid id")
	}
	user, err := s.queries.GetUserByID(c, pgtype.UUID{Valid: true, Bytes: userID})
	if errors.Is(err, sql.ErrNoRows) {
		user, err = s.fetchKeycloakUser(c, userID)
		if err != nil {
			return PublishResponseError(100, "Internal server error")
		}
	} else if err != nil {
		return PublishResponseError(100, "Internal server error")
	}
	res, err := s.queries.UserCanPublish(c, sqlc.UserCanPublishParams{
		ID:      user.ID,
		Channel: req.Channel,
	})

	if err != nil {
		return PublishResponseError(100, "Internal server error")
	}

	if res == 0 {
		return PublishResponseError(103, "permission denied")
	}
	return &proxyproto.PublishResponse{}, nil
}

func (s *Service) GetCloakByUser(c context.Context, userId string) (*gocloak.User, error) {
	log.Println("GetCloakByUser")
	if s.token == nil || s.expiredAt.After(time.Now()) {
		token, err := s.kcClient.LoginClient(c, s.KeyCloakClient, s.KeyCloakSecret, s.KeyCloakRealm)
		if err != nil {
			return nil, err
		}
		s.token = token
		s.expiredAt = time.Now().Add(time.Second * time.Duration(s.token.ExpiresIn))
	}
	user, err := s.kcClient.GetUserByID(c, s.token.AccessToken, s.KeyCloakRealm, userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
