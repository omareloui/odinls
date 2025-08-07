package client

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ClientService interface {
	GetClients(claims *jwtadapter.AccessClaims) ([]Client, error)
	GetClientByID(claims *jwtadapter.AccessClaims, id string) (*Client, error)
	CreateClient(claims *jwtadapter.AccessClaims, client *Client) (*Client, error)
	UpdateClientByID(claims *jwtadapter.AccessClaims, id string, client *Client) (*Client, error)
}
