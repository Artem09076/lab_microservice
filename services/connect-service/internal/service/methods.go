package service

import (
	"context"
	"encoding/json"
	"s0709-22/internal/proxyproto"
	"strconv"
)

// Connect ...
func (s *Service) Connect(ctx context.Context, request *proxyproto.ConnectRequest) (*proxyproto.ConnectResponse, error) {
	type AuthRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	authRequest := &AuthRequest{}

	if err := json.Unmarshal(request.Data, authRequest); err != nil {
		return RespondError(107, "bad request")
	}

	account, err := s.storage.GetUserByUsermame(ctx, authRequest.Username)
	if err != nil {
		return RespondError(101, "unauthorized")
	}

	if authRequest.Password != account.Password {
		return RespondError(101, "unauthorized")
	}

	return &proxyproto.ConnectResponse{
		Result: &proxyproto.ConnectResult{
			User: strconv.FormatInt(account.ID, 10),
		},
	}, nil
}
