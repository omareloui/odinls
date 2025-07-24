package material

import jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"

type MaterialService interface {
	GetMaterials(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Material, error)
	GetMaterialByID(claims *jwtadapter.JwtAccessClaims, id string, opts ...RetrieveOptsFunc) (*Material, error)
	CreateMaterial(claims *jwtadapter.JwtAccessClaims, mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
	UpdateMaterialByID(claims *jwtadapter.JwtAccessClaims, id string, mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
}
