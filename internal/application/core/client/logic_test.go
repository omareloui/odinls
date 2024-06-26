package client_test

import (
	"testing"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/client"
	client_mock "github.com/omareloui/odinls/internal/application/core/client/mocks"
	"github.com/omareloui/odinls/internal/application/core/role"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/sanitizer/conformadaptor"
	"github.com/omareloui/odinls/internal/validator/playgroundvalidator"
	"github.com/stretchr/testify/assert"
)

func TestGetClients(t *testing.T) {
	mockRepo := new(client_mock.MockClientRepository)

	clients := []client.Client{}
	mockRepo.On("GetClients").Return(clients, nil)

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := client.NewClientService(mockRepo, v, sani)

	t.Run("with permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{
				Name: role.OPAdmin.String(),
			},
			CraftsmanInfo: user.Craftsman{
				MerchantID: "1234",
			},
		}

		actualClients, err := s.GetClients(&claims)
		mockRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, clients, actualClients)
	})

	t.Run("without permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{
				Name: role.Admin.String(),
			},
		}

		actualClients, err := s.GetClients(&claims)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})

	t.Run("without claims", func(t *testing.T) {
		actualClients, err := s.GetClients(nil)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})
}

func TestGetCurrentMerchantClients(t *testing.T) {
	mockRepo := new(client_mock.MockClientRepository)

	merId := "1234"

	clients := []client.Client{{MerchantID: merId}}
	mockRepo.On("GetClientsByMerchantID", merId).Return(clients, nil)

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := client.NewClientService(mockRepo, v, sani)

	t.Run("is craftsman and with permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{
				Name: role.Moderator.String(),
			},
			CraftsmanInfo: user.Craftsman{
				MerchantID: merId,
			},
		}

		actualClients, err := s.GetCurrentMerchantClients(&claims)
		mockRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, clients, actualClients)
	})

	t.Run("is craftsman and without permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{
				Name: role.NoAuthority.String(),
			},
			CraftsmanInfo: user.Craftsman{
				MerchantID: merId,
			},
		}

		actualClients, err := s.GetCurrentMerchantClients(&claims)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})

	t.Run("not craftsman and with permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{
				Name: role.NoAuthority.String(),
			},
		}

		actualClients, err := s.GetCurrentMerchantClients(&claims)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})

	t.Run("no claims", func(t *testing.T) {
		actualClients, err := s.GetCurrentMerchantClients(nil)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})
}

func TestGetClientByID(t *testing.T) {
	mockRepo := new(client_mock.MockClientRepository)

	clientId := "11"
	merId := "1234"

	cli := client.Client{ID: clientId, MerchantID: merId}
	mockRepo.On("GetClientByID", clientId).Return(&cli, nil)

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := client.NewClientService(mockRepo, v, sani)

	t.Run("is craftsman", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			CraftsmanInfo: user.Craftsman{
				MerchantID: merId,
			},
		}

		actualClient, err := s.GetClientByID(&claims, clientId)
		mockRepo.AssertExpectations(t)
		assert.Nil(t, err)
		assert.Equal(t, &cli, actualClient)
	})

	t.Run("not craftsman", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{}

		actualClients, err := s.GetClientByID(&claims, clientId)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})

	t.Run("no claims", func(t *testing.T) {
		actualClients, err := s.GetClientByID(nil, clientId)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
		assert.Nil(t, actualClients)
	})
}

func TestCreateClient(t *testing.T) {
	mockRepo := new(client_mock.MockClientRepository)

	clientId := "11"
	merId := "1234"

	cli := client.Client{
		ID:    clientId,
		Name:  "mock client name",
		Notes: "",
		ContactInfo: client.ContactInfo{
			PhoneNumbers: map[string]string{},
			Emails:       map[string]string{},
			Links:        map[string]string{},
			Locations:    map[string]string{},
		},
		WholesaleAsDefault: false,
	}

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := client.NewClientService(mockRepo, v, sani)

	t.Run("permissions", func(t *testing.T) {
		cli2 := cli
		mockRepo.On("CreateClient", &cli2).Return(nil)

		t.Run("with permissions", func(t *testing.T) {
			claims := jwtadapter.JwtAccessClaims{
				Role:          role.Role{Name: role.Admin.String()},
				CraftsmanInfo: user.Craftsman{MerchantID: merId},
			}

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, merId, cli2.MerchantID)
			assert.Equal(t, clientId, cli2.ID)
		})

		t.Run("without permissions", func(t *testing.T) {
			claims := jwtadapter.JwtAccessClaims{
				Role: role.Role{Name: role.Moderator.String()},
			}

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)
			assert.ErrorIs(t, errs.ErrForbidden, err)
		})

		t.Run("no claims", func(t *testing.T) {
			err := s.CreateClient(nil, &cli2)
			mockRepo.AssertExpectations(t)
			assert.ErrorIs(t, errs.ErrForbidden, err)
		})
	})

	t.Run("validation and sanitization", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role:          role.Role{Name: role.Admin.String()},
			CraftsmanInfo: user.Craftsman{MerchantID: merId},
		}

		t.Run("valid inputs", func(t *testing.T) {
			cli2 := cli
			mockRepo.On("CreateClient", &cli2).Return(nil)
			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, merId, cli2.MerchantID)
			assert.Equal(t, clientId, cli2.ID)
		})

		t.Run("nil contact info maps", func(t *testing.T) {
			cli2 := cli
			mockRepo.On("CreateClient", &cli2).Return(nil)
			cli2.ContactInfo = client.ContactInfo{}

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, merId, cli2.MerchantID)
			assert.Equal(t, clientId, cli2.ID)
		})

		t.Run("no name", func(t *testing.T) {
			cli2 := cli
			cli2.Name = ""

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.NotNil(t, err)
			assert.Equal(t, errs.ValidationError{}.Error(), err.Error())
		})

		t.Run("validate contact info keys and values", func(t *testing.T) {
			cli2 := cli
			cli2.ContactInfo = client.ContactInfo{
				Links:  map[string]string{"fb": "nohting", "ig ": "https://ig.com"},
				Emails: map[string]string{"default": "invalidemail"},
				Locations: map[string]string{
					"home": "             ",
					"    ": "va",
				},
				PhoneNumbers: map[string]string{
					"key":  "val",
					"key2": "01111",
				},
			}

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.NotNil(t, err)
			assert.Equal(t, errs.ValidationError{}.Error(), err.Error())

			if valerr, ok := err.(errs.ValidationError); ok {
				assert.Equal(t, valerr.Errors["ContactInfo.Emails[default]"].Msg(), "Invalid email")
				assert.Contains(t, valerr.Errors["ContactInfo.Links[ig]"].Msg(), "Value is too short")
				assert.Contains(t, valerr.Errors["ContactInfo.Links[fb]"].Msg(), "valid URL")
				assert.Contains(t, valerr.Errors["ContactInfo.Locations[home]"].Msg(), "required")
				assert.Contains(t, valerr.Errors["ContactInfo.Locations[]"].Msg(), "required")
			}
		})

		t.Run("sanitize contact info keys and values", func(t *testing.T) {
			cli2 := cli
			mockRepo.On("CreateClient", &cli2).Return(nil)
			notSanitizedPhoneKey := " 		home"
			notSanitizedPhoneValue := " 01111  "
			notSanitizedLinkKey := " 		facebook "
			notSanitizedLinkValue := " https://fb.com  	"
			sanitizedPhoneKey := "home"
			sanitizedPhoneValue := "01111"
			sanitizedLinkKey := "facebook"
			sanitizedLinkValue := "https://fb.com"
			cli2.ContactInfo = client.ContactInfo{
				PhoneNumbers: map[string]string{notSanitizedPhoneKey: notSanitizedPhoneValue},
				Links:        map[string]string{notSanitizedLinkKey: notSanitizedLinkValue},
			}

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)

			_, ok := cli2.ContactInfo.PhoneNumbers[notSanitizedPhoneKey]
			assert.False(t, ok)
			_, ok = cli2.ContactInfo.Links[notSanitizedLinkKey]
			assert.False(t, ok)

			assert.Equal(t, cli2.ContactInfo.PhoneNumbers[sanitizedPhoneKey], sanitizedPhoneValue)
			assert.Equal(t, cli2.ContactInfo.Links[sanitizedLinkKey], sanitizedLinkValue)
		})

		t.Run("trim spaces and tabs", func(t *testing.T) {
			cli2 := cli
			mockRepo.On("CreateClient", &cli2).Return(nil)

			notSanitizedName := " 		seif eloui "
			sanitizedName := "Seif Eloui"
			notSanitizedNotes := " 		    this is a note!          "
			sanitizedNotes := "this is a note!"

			cli2.Name = notSanitizedName
			cli2.Notes = notSanitizedNotes

			err := s.CreateClient(&claims, &cli2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, sanitizedName, cli2.Name)
			assert.Equal(t, sanitizedNotes, cli2.Notes)
		})
	})
}

func TestUpdateClient(t *testing.T) {
	mockRepo := new(client_mock.MockClientRepository)

	clientId := "11"
	merId := "1234"

	cli := client.Client{
		ID:    clientId,
		Name:  "mock client name",
		Notes: "nothing",
		ContactInfo: client.ContactInfo{
			PhoneNumbers: map[string]string{},
			Emails:       map[string]string{},
			Links:        map[string]string{},
			Locations:    map[string]string{},
		},
		WholesaleAsDefault: false,
	}

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := client.NewClientService(mockRepo, v, sani)

	t.Run("with permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role:          role.Role{Name: role.Admin.String()},
			CraftsmanInfo: user.Craftsman{MerchantID: merId},
		}

		cli2 := cli
		mockRepo.On("UpdateClientByID", clientId, &cli2).Return(nil)

		cli2.CreatedAt = time.Now()

		err := s.UpdateClientByID(&claims, clientId, &cli2)
		mockRepo.AssertExpectations(t)

		assert.Nil(t, err)
		assert.Equal(t, clientId, cli2.ID)
		assert.True(t, cli2.CreatedAt.IsZero())
	})

	t.Run("without permissions", func(t *testing.T) {
		claims := jwtadapter.JwtAccessClaims{
			Role: role.Role{Name: role.Moderator.String()},
		}

		cli2 := cli
		err := s.UpdateClientByID(&claims, clientId, &cli2)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
	})

	t.Run("no claims", func(t *testing.T) {
		cli2 := cli

		err := s.UpdateClientByID(nil, clientId, &cli2)
		mockRepo.AssertExpectations(t)
		assert.ErrorIs(t, errs.ErrForbidden, err)
	})
}
