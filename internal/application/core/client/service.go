package client

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type ClientService interface {
	GetClients(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Client, error)
	GetCurrentMerchantClients(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Client, error)
	GetClientByID(claims *jwtadapter.JwtAccessClaims, id string, opts ...RetrieveOptsFunc) (*Client, error)
	CreateClient(claims *jwtadapter.JwtAccessClaims, client *Client, opts ...RetrieveOptsFunc) error
	UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, client *Client, opts ...RetrieveOptsFunc) error
}
