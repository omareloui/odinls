package client

import (
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/interfaces"
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

func (s *clientService) GetClients(claims *jwtadapter.JwtAccessClaims) ([]Client, error) {
	if claims == nil || (!claims.Role.IsModerator() && !claims.IsCraftsman()) {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClients()
}

func (s *clientService) GetClientByID(claims *jwtadapter.JwtAccessClaims, id string) (*Client, error) {
	if claims == nil || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	return s.repo.GetClientByID(id)
}

func (s *clientService) CreateClient(claims *jwtadapter.JwtAccessClaims, client *Client) (*Client, error) {
	if claims == nil || !claims.Role.IsAdmin() || !claims.IsCraftsman() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizeClient(client)
	if err != nil {
		return nil, err
	}

	if err := s.validator.Validate(client); err != nil {
		return nil, s.validator.ParseError(err)
	}

	return s.repo.CreateClient(client)
}

func (s *clientService) UpdateClientByID(claims *jwtadapter.JwtAccessClaims, id string, client *Client) (*Client, error) {
	if claims == nil || !claims.Role.IsAdmin() {
		return nil, errs.ErrForbidden
	}

	err := s.sanitizeClient(client)
	if err != nil {
		return nil, err
	}

	if err := s.validator.Validate(client); err != nil {
		return nil, s.validator.ParseError(err)
	}

	return s.repo.UpdateClientByID(id, client)
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
