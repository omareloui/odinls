package material

import (
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

type materialService struct {
	repo      MaterialRepository
	validator interfaces.Validator
	sanitizer interfaces.Sanitizer
}

func NewMaterialService(repo MaterialRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *materialService {
	return &materialService{
		repo:      repo,
		validator: validator,
		sanitizer: sanitizer,
	}
}

func (s *materialService) GetMaterials(claims *jwtadapter.JwtAccessClaims, options ...RetrieveOptsFunc) ([]Material, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetMaterials(options...)
}

func (s *materialService) GetMaterialByID(claims *jwtadapter.JwtAccessClaims, id string, options ...RetrieveOptsFunc) (*Material, error) {
	if claims == nil || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetMaterialByID(id, options...)
}

func (s *materialService) CreateMaterial(claims *jwtadapter.JwtAccessClaims, mat *Material, options ...RetrieveOptsFunc) (*Material, error) {
	if claims == nil || !claims.Role.IsAdmin() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(mat)
	if err != nil {
		return nil, errs.ErrSanitizer
	}

	if err := s.validator.Validate(mat); err != nil {
		return nil, s.validator.ParseError(err)
	}

	return s.repo.CreateMaterial(mat, options...)
}

func (s *materialService) UpdateMaterialByID(claims *jwtadapter.JwtAccessClaims, id string, umat *Material, options ...RetrieveOptsFunc) (*Material, error) {
	if claims == nil || !claims.Role.IsAdmin() {
		return errs.ErrForbidden
	}

	err := s.sanitizer.SanitizeStruct(umat)
	if err != nil {
		return errs.ErrSanitizer
	}

	if err := s.validator.Validate(umat); err != nil {
		return s.validator.ParseError(err)
	}

	mat, err := s.repo.GetMaterialByID(id)
	if err != nil {
		return err
	}

	umat.CreatedAt = mat.CreatedAt
	umat.UpdatedAt = time.Now()

	return s.repo.UpdateMaterialByID(id, umat, options...)
}
