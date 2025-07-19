package client

import (
	"errors"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
)

var (
	ErrClientNotFound = errors.New("client not found")
	ErrClientExists   = errors.New("client exists")
)

type clientService struct {
	repo      ClientRepository
	validator interfaces.Validator
	sanitizer interfaces.Sanitizer
}

func NewClientService(clientRepository ClientRepository, validator interfaces.Validator, sanitizer interfaces.Sanitizer) *clientService {
	return &clientService{
		repo:      clientRepository,
		validator: validator,
		sanitizer: sanitizer,
	}
}

func (s *clientService) GetClients(claims *jwtadapter.JwtAccessClaims, opts ...RetrieveOptsFunc) ([]Client, error) {
	if claims == nil || (!claims.Role.IsModerator() && !claims.IsCraftsman()) {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClients(opts...)
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

	err := s.sanitizeClient(client)
	if err != nil {
		return err
	}

	if err := s.validator.Validate(client); err != nil {
		return s.validator.ParseError(err)
	}

	now := time.Now()
	client.CreatedAt = now
	client.UpdatedAt = now

	return s.repo.CreateClient(client, opts...)
}

func (s *clientService) UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, client *Client, opts ...RetrieveOptsFunc) error {
	if claims == nil || !claims.Role.IsAdmin() {
		return errs.ErrForbidden
	}

	err := s.sanitizeClient(client)
	if err != nil {
		return err
	}

	if err := s.validator.Validate(client); err != nil {
		return s.validator.ParseError(err)
	}

	// TODO(refactor): make sure to update the updated from the SERVICE level in all services
	// TODO(refactor): DON'T as this is a database responsibility, move everything there.

	client.CreatedAt = time.Time{}
	client.UpdatedAt = time.Now()

	return s.repo.UpdateClientByID(id, client, opts...)
}

func (s *clientService) sanitizeClient(cli *Client) error {
	err := s.sanitizer.SanitizeStruct(cli)
	if err != nil {
		return errs.ErrSanitizer
	}

	cli.ContactInfo.Locations = s.sanitizer.TrimMap(cli.ContactInfo.Locations)
	cli.ContactInfo.Emails = s.sanitizer.TrimMap(cli.ContactInfo.Emails)
	cli.ContactInfo.Links = s.sanitizer.TrimMap(cli.ContactInfo.Links)
	cli.ContactInfo.PhoneNumbers = s.sanitizer.TrimMap(cli.ContactInfo.PhoneNumbers)
	return nil
}
