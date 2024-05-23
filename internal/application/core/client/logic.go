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

	return s.repo.GetClientByID(claims.CraftsmanInfo.MerchantID, opts...)
}

func (s *clientService) CreateClient(claims *jwtadapter.JwtAccessClaims, client *Client, opts ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() {
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

	// TODO(refactor): change all patch request to put
	// TODO(refactor): make sure to update the updated from the SERVICE level in all services

	client.CreatedAt = time.Time{}
	client.UpdatedAt = time.Now()

	return s.repo.UpdateClientByID(id, client, opts...)
}

func sanitizeClient(c *Client) {
	c.Name = sanitizer.TrimString(c.Name)
	c.Notes = sanitizer.TrimString(c.Notes)
	sanitizeMap(&c.ContactInfo.Emails)
	sanitizeMap(&c.ContactInfo.Links)
	sanitizeMap(&c.ContactInfo.Locations)
	sanitizeMap(&c.ContactInfo.PhoneNumber)
}

func sanitizeMap(m *map[string]string) {
	if m == nil {
		return
	}
	old := *m
	newm := map[string]string{}
	for key, val := range old {
		newm[sanitizer.TrimString(key)] = sanitizer.TrimString(val)
	}
	*m = newm
}
