package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Artem09076/lab_microservice.git/internal/proxyproto"
	sqlc "github.com/Artem09076/lab_microservice.git/internal/userdb"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) fetchKeycloakUser(c context.Context, userId uuid.UUID) (sqlc.User, error) {
	userKC, err := s.kcClient.GetUserByID(c, userId.String())
	if err != nil {
		return sqlc.User{}, err
	}
	user := sqlc.User{
		ID:         pgtype.UUID{Valid: true, Bytes: userId},
		Username:   userKC.Username,
		GivenName:  userKC.FirstName,
		FamilyName: userKC.LastName,
		Enabled:    userKC.Enabled,
	}

	if err = s.queries.CreateUser(c, sqlc.CreateUserParams{
		ID:         pgtype.UUID{Valid: true, Bytes: userId},
		Username:   userKC.Username,
		GivenName:  userKC.FirstName,
		FamilyName: userKC.LastName,
		Enabled:    userKC.Enabled,
	}); err != nil {
		return sqlc.User{}, err
	}
	return user, nil
}

func (s *Service) Subscribe(c context.Context, req *proxyproto.SubscribeRequest) (*proxyproto.SubscribeResponse, error) {
	userID, err := uuid.Parse(req.User)
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
