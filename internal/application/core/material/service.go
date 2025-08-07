package material

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type MaterialService interface {
	GetMaterials(claims *jwtadapter.AccessClaims, opts ...RetrieveOptsFunc) ([]Material, error)
	GetMaterialByID(claims *jwtadapter.AccessClaims, id string, opts ...RetrieveOptsFunc) (*Material, error)
	CreateMaterial(claims *jwtadapter.AccessClaims, mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
	UpdateMaterialByID(claims *jwtadapter.AccessClaims, id string, mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
}
