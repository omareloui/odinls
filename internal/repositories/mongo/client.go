package mongo

import (
	"github.com/omareloui/odinls/internal/application/core/client"
)

func (r *repository) GetClients() ([]client.Client, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetAll[client.Client](ctx, r.clientsColl)
}

func (r *repository) GetClientByID(id string) (*client.Client, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return GetByID[client.Client](ctx, r.clientsColl, id)
}

func (r *repository) CreateClient(cli *client.Client) (*client.Client, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return InsertStruct(ctx, r.clientsColl, cli)
}

func (r *repository) UpdateClientByID(id string, cli *client.Client) (*client.Client, error) {
	ctx, cancel := r.newCtx()
	defer cancel()

	return UpdateStructByID[client.Client](ctx, r.clientsColl, id, cli)
}
