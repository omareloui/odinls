package client

type ClientRepository interface {
	GetClients() ([]Client, error)
	GetClientByID(id string) (*Client, error)
	CreateClient(client *Client) (*Client, error)
	UpdateClientByID(id string, client *Client) (*Client, error)
}
