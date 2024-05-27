package client

import (
	"errors"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
	"github.com/omareloui/odinls/internal/sanitizer"
)

var (
	ErrClientNotFound          = errors.New("client not found")
	ErrClientExistsForMerchant = errors.New("client exists for that merchant")
)

type clientService struct {
	repo      ClientRepository
	validator interfaces.Validator
}

func NewClientService(clientRepository ClientRepository, validator interfaces.Validator) *clientService {
	return &clientService{
		repo:      clientRepository,
		validator: validator,
	}
}

func (s *clientService) GetClients(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Client, error) {
	if claims == nil || !claims.Role.IsOPAdmin() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClients(opts...)
}

func (s *clientService) GetCurrentMerchantClients(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Client, error) {
	if claims == nil || !claims.IsCraftsman() || !claims.Role.IsModerator() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClientsByMerchantID(claims.CraftsmanInfo.MerchantID, opts...)
}

func (s *clientService) GetClientByID(claims *jwtadapter.JwtAccessClaims, id string, opts ...RetrieveOptsFunc) (*Client, error) {
	if claims == nil || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClientByID(id, opts...)
}

func (s *clientService) CreateClient(claims *jwtadapter.JwtAccessClaims, client *Client, opts ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return errs.ErrForbidden
	}

	sanitizeClient(client)

	if err := s.validator.Validate(client); err != nil {
		return s.validator.ParseError(err)
	}

	client.MerchantID = claims.CraftsmanInfo.MerchantID

	now := time.Now()
	client.CreatedAt = now
	client.UpdatedAt = now

	return s.repo.CreateClient(client, opts...)
}

func (s *clientService) UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, client *Client, opts ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() {
		return errs.ErrForbidden
	}

	sanitizeClient(client)

	if err := s.validator.Validate(client); err != nil {
		return s.validator.ParseError(err)
	}

	// TODO(refactor): make sure to update the updated from the SERVICE level in all services
	// TODO(refactor): DON'T as this is a database responsiblity, move everything there.

	client.CreatedAt = time.Time{}
	client.UpdatedAt = time.Now()

	return s.repo.UpdateClientByID(id, client, opts...)
}

func sanitizeClient(c *Client) {
	c.Name = sanitizer.TrimString(c.Name)
	c.Notes = sanitizer.TrimString(c.Notes)
	sanitizer.SanitizeStringMap(&c.ContactInfo.Emails)
	sanitizer.SanitizeStringMap(&c.ContactInfo.Links)
	sanitizer.SanitizeStringMap(&c.ContactInfo.Locations)
	sanitizer.SanitizeStringMap(&c.ContactInfo.PhoneNumbers)
}
