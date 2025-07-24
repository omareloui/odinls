package material

type MaterialRepository interface {
	GetMaterials(opts ...RetrieveOptsFunc) ([]Material, error)
	GetMaterialByID(id string, opts ...RetrieveOptsFunc) (*Material, error)
	CreateMaterial(mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
	UpdateMaterialByID(id string, mat *Material, opts ...RetrieveOptsFunc) (*Material, error)
}
