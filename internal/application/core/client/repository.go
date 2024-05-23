package client

type ClientRepository interface {
	GetClients(opts ...RetrieveOptsFunc) ([]Client, error)
	GetClientsByMerchantID(merchantId string, opts ...RetrieveOptsFunc) ([]Client, error)
	GetClientByID(id string, opts ...RetrieveOptsFunc) (*Client, error)
	CreateClient(client *Client, opts ...RetrieveOptsFunc) error
	UpdateClientByID(id string, client *Client, opts ...RetrieveOptsFunc) error
}
