package kc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type accessToken struct {
	AccessToken string
	ExpiresAt   time.Time
}

type KCClient struct {
	baseUrl      string
	realmName    string
	clientId     string
	clientSecret string
	token        *accessToken
}

type KCUser struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Enabled   bool   `json:"enabled"`
}

func New(baseUrl, realmName, clientID, clientSecret string) *KCClient {
	return &KCClient{
		baseUrl:      baseUrl,
		realmName:    realmName,
		clientId:     clientID,
		clientSecret: clientSecret,
	}
}

func (c *KCClient) tokenRefresh(ctx context.Context) error {
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.baseUrl, c.realmName)
	data := url.Values{}
	data.Add("grant_type", "client_credentials")
	data.Add("client_id", c.clientId)
	data.Add("client_secret", c.clientSecret)
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		tokenURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("keycloak: bad request")
	}
	type tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
	}
	token := tokenResponse{}

	tokenDec := json.NewDecoder(response.Body)

	if err := tokenDec.Decode(token); err != nil {
		return err
	}
	c.token = &accessToken{
		AccessToken: token.AccessToken,
		ExpiresAt:   time.Now().Add(time.Duration(token.ExpiresIn) * time.Second),
	}
	return nil
}

func (c *KCClient) sendRequestWithToken(ctx context.Context, request *http.Request) (*http.Response, error) {
	if c.token == nil || c.token.ExpiresAt.Before(time.Now()) {
		err := c.tokenRefresh(ctx)
		if err != nil {
			return nil, err
		}
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token.AccessToken))

	body, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *KCClient) GetUserByID(ctx context.Context, userId string) (*KCUser, error) {
	userUrl := fmt.Sprintf(
		"%s/admin/realms/%s/users/%s",
		c.baseUrl,
		c.realmName,
		userId,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", userUrl, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.sendRequestWithToken(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("keycloak: user not found")
	} else if resp.StatusCode != 200 {
		return nil, errors.New("keycloak: bad request")
	}

	user := KCUser{}
	userDecoder := json.NewDecoder(resp.Body)
	err = userDecoder.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
