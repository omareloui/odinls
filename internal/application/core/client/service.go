package client

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ClientService interface {
	GetClients(claims *jwtadapter.JwtAccessClaims) ([]Client, error)
	GetClientByID(claims *jwtadapter.JwtAccessClaims, id string) (*Client, error)
	CreateClient(claims *jwtadapter.JwtAccessClaims, client *Client) (*Client, error)
	UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, client *Client) (*Client, error)
}
