package order_test

import (
	"testing"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	counter_mock "github.com/omareloui/odinls/internal/application/core/counter/mocks"
	"github.com/omareloui/odinls/internal/application/core/order"
	order_mock "github.com/omareloui/odinls/internal/application/core/order/mocks"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/application/core/user"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/internal/sanitizer/conformadaptor"
	"github.com/omareloui/odinls/internal/validator/playgroundvalidator"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	mockProdS := new(product_mock.MockProductService)
	mockCounterS := new(counter_mock.MockCounterService)
	mockRepo := new(order_mock.MockOrderRepository)

	orderId := "665dbe5ac352603c7e73da4f"

	prod := product.Product{
		ID: "665dbe5ac352603c7e68fa5e",
		Variants: []product.Variant{{
			ID:    "665dbe5ac352610c7e73fa5e",
			Price: 300,
		}},
	}

	ord := order.Order{
		ID:       orderId,
		ClientID: "665dbe5ac352603c7e73fa5e",
		Status:   order.StatusPendingConfirmation.String(),
		Timeline: order.Timeline{
			IssuanceDate: time.Now(),
		},
		Items: []order.Item{
			{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
		},
	}

	v := playgroundvalidator.NewValidator()
	sani := conformadaptor.NewSanitizer()
	s := order.NewOrderService(mockRepo, mockProdS, mockCounterS, v, sani)

	t.Run("permissions", func(t *testing.T) {
		ord2 := ord
		mockRepo.On("CreateOrder", &ord2).Return(nil)

		t.Run("with permissions", func(t *testing.T) {
			claims := jwtadapter.JwtAccessClaims{
				Role:          user.Admin,
				CraftsmanInfo: &user.Craftsman{},
			}
			mockProdS.On("GetProductByIDAndVariantID", &claims, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(uint(55), nil)

			err := s.CreateOrder(&claims, &ord2)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, orderId, ord2.ID)
		})

		t.Run("without permissions", func(t *testing.T) {
			claims := jwtadapter.JwtAccessClaims{
				Role: user.Moderator,
			}

			err := s.CreateOrder(&claims, &ord2)
			mockRepo.AssertExpectations(t)
			assert.ErrorIs(t, errs.ErrForbidden, err)
		})

		t.Run("no claims", func(t *testing.T) {
			err := s.CreateOrder(nil, &ord2)
			mockRepo.AssertExpectations(t)
			assert.ErrorIs(t, errs.ErrForbidden, err)
		})
	})

	t.Run("validation", func(t *testing.T) {
		orderId := "665dbe5ac352603c7e73da4f"

		prod := product.Product{
			ID: "665dbe5ac352603c7e68fa5e",
			Variants: []product.Variant{{
				ID:    "665dbe5ac352610c7e73fa5e",
				Price: 300,
			}},
		}

		ord2 := order.Order{
			ID:       orderId,
			ClientID: "665dbe5ac352603c7e73fa5e",
			Status:   order.StatusPendingConfirmation.String(),
			Timeline: order.Timeline{
				IssuanceDate: time.Now(),
			},
			Items: []order.Item{
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
			},
		}

		claims := jwtadapter.JwtAccessClaims{
			Role:          role.Role{Name: role.SuperAdmin.String()},
			CraftsmanInfo: user.Craftsman{},
		}

		t.Run("valid inputs", func(t *testing.T) {
			ordNum := uint(18356)
			ord3 := ord2
			mockRepo.On("CreateOrder", &ord3).Return(nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(ordNum, nil)

			err := s.CreateOrder(&claims, &ord3)
			mockRepo.AssertExpectations(t)

			assert.Nil(t, err)
			assert.Equal(t, orderId, ord3.ID)
			assert.Equal(t, ordNum, ord3.Number)
		})

		t.Run("items", func(t *testing.T) {
			t.Run("invalid item array length", func(t *testing.T) {
				ord3 := ord2
				ord3.Items = []order.Item{}

				err := s.CreateOrder(&claims, &ord3)

				assert.NotNil(t, err)
				if valerr, ok := err.(errs.ValidationError); ok {
					assert.Contains(t, valerr.Errors["Items"].Msg(), "too short")
				}
			})

			t.Run("invalid item product id", func(t *testing.T) {
				ord3 := ord2

				prodIdExsists := ord3.Items[0].ProductID
				prodIdDoesntExist := "665dbe5ac352603c7e68fa1e"

				mockProdS.On("GetProductByIDAndVariantID", &claims, prodIdDoesntExist, ord3.Items[0].VariantID).Return(nil, product.ErrProductNotFound)
				mockProdS.On("GetProductByIDAndVariantID", &claims, ord3.Items[0].ProductID, ord3.Items[0].VariantID).Return(&prod, nil)
				mockRepo.On("CreateOrder", &ord3).Return(nil)

				ord3.Items[0].ProductID = prodIdDoesntExist
				err := s.CreateOrder(&claims, &ord3)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, product.ErrProductNotFound)

				ord3.Items[0].ProductID = prodIdExsists
				err = s.CreateOrder(&claims, &ord3)
				assert.Nil(t, err)
			})

			t.Run("invalid item variant id", func(t *testing.T) {
				ord3 := ord2

				variantIdExsists := ord3.Items[0].VariantID
				variantIdDoesntExist := "665dbe5ac352603c7e68fa1e"

				mockProdS.On("GetProductByIDAndVariantID", &claims, ord3.Items[0].ProductID, variantIdDoesntExist).Return(nil, product.ErrProductNotFound)
				mockProdS.On("GetProductByIDAndVariantID", &claims, ord3.Items[0].ProductID, variantIdExsists).Return(&prod, nil)
				mockRepo.On("CreateOrder", &ord3).Return(nil)

				ord3.Items[0].VariantID = variantIdDoesntExist
				err := s.CreateOrder(&claims, &ord3)
				assert.NotNil(t, err)
				assert.ErrorIs(t, err, product.ErrProductNotFound)

				ord3.Items[0].VariantID = variantIdExsists
				err = s.CreateOrder(&claims, &ord3)
				assert.Nil(t, err)
			})
		})

		t.Run("price addon", func(t *testing.T) {
			t.Run("invalid price addon", func(t *testing.T) {
				ord3 := ord2
				ord3.PriceAddons = []order.PriceAddon{{
					Kind:         "invalid_kind",
					Amount:       -1,
					IsPercentage: false,
				}}

				err := s.CreateOrder(&claims, &ord3)
				assert.NotNil(t, err)
				if valerr, ok := err.(errs.ValidationError); ok {
					assert.Contains(t, valerr.Errors["PriceAddons[0].Amount"].Msg(), "Value is low")
					assert.Contains(t, valerr.Errors["PriceAddons[0].Kind"].Msg(), "one of")
				}
			})
		})
	})

	t.Run("filling and population", func(t *testing.T) {
		orderId := "665dbe5ac352603c7e73da4f"

		prod := product.Product{
			ID: "665dbe5ac352603c7e68fa5e",
			Variants: []product.Variant{{
				ID:    "665dbe5ac352610c7e73fa5e",
				Price: 300,
			}},
		}

		ord := order.Order{
			ID:       orderId,
			ClientID: "665dbe5ac352603c7e73fa5e",
			Status:   order.StatusPendingConfirmation.String(),
			Timeline: order.Timeline{
				IssuanceDate: time.Now(),
			},
			Items: []order.Item{
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
			},
		}

		claims := jwtadapter.JwtAccessClaims{
			Role:          user.Admin,
			CraftsmanInfo: &user.Craftsman{},
		}

		t.Run("fills the order number", func(t *testing.T) {
			claims2 := claims

			ord2 := ord

			num := uint(100)

			mockRepo.On("CreateOrder", &ord2).Return(nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims2, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims2).Return(num, nil)

			err := s.CreateOrder(&claims2, &ord2)

			assert.Nil(t, err)
			assert.Equal(t, num, ord2.Number)
		})

		t.Run("generate a new ref", func(t *testing.T) {
			num := uint(14)
			ord2 := ord
			mockRepo.On("CreateOrder", &ord2).Return(nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(num, nil)

			err := s.CreateOrder(&claims, &ord2)

			assert.Nil(t, err)
			assert.NotZero(t, ord2.Ref)
			assert.Len(t, ord2.Ref, 8)
		})

		t.Run("fills order item price", func(t *testing.T) {
			claims2 := claims
			ord2 := ord
			ord2.Items[0].Price = float64(45.99)
			mockRepo.On("CreateOrder", &ord2).Return(nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims2, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims2).Return(uint(10), nil)

			err := s.CreateOrder(&claims2, &ord2)

			assert.Nil(t, err)
			assert.Equal(t, prod.Variants[0].Price, ord2.Items[0].Price)
		})

		t.Run("fills the order item default progress", func(t *testing.T) {
			claims2 := claims
			ord2 := ord
			ord2.Items[0].Price = float64(45.99)
			mockRepo.On("CreateOrder", &ord2).Return(nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims2, ord2.Items[0].ProductID, ord2.Items[0].VariantID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims2).Return(uint(10), nil)

			err := s.CreateOrder(&claims2, &ord2)

			assert.Nil(t, err)
			assert.Equal(t, order.ItemProgressNotStarted.String(), ord2.Items[0].Progress)
		})
	})

	t.Run("calculations", func(t *testing.T) {
		orderId := "665dbe5ac352603c7e73da4b"

		prod := product.Product{
			ID: "665dbe5ac352603c7e68fa5a",
			Variants: []product.Variant{{
				ID:    "665dbe5ac352610c7e73fa5f",
				Price: 300,
			}, {
				ID:    "665dbe5ac352610c7e73fa4f",
				Price: 200,
			}},
		}

		ord := order.Order{
			ID:       orderId,
			ClientID: "665dbe5ac352603c7e73fa5c",
			Status:   order.StatusPendingConfirmation.String(),
			Timeline: order.Timeline{
				IssuanceDate: time.Now(),
			},
			Items: []order.Item{
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
				{ProductID: prod.ID, VariantID: prod.Variants[1].ID},
			},
		}

		claims := jwtadapter.JwtAccessClaims{
			Role:          user.Admin,
			CraftsmanInfo: &user.Craftsman{},
		}

		t.Run("sum items prices without price addons", func(t *testing.T) {
			ord2 := ord

			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[0].ID).Return(&prod, nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[1].ID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(uint(32), nil)

			mockRepo.On("CreateOrder", &ord2).Return(nil)

			err := s.CreateOrder(&claims, &ord2)
			assert.Nil(t, err)

			assert.Equal(t, float64(800), ord2.Subtotal)
		})

		t.Run("calculate subtotal with items and absolute addons", func(t *testing.T) {
			ord2 := ord

			ord2.PriceAddons = append(ord2.PriceAddons, order.PriceAddon{
				Kind:         "fees",
				Amount:       50,
				IsPercentage: false,
			}, order.PriceAddon{
				Kind:         "taxes",
				Amount:       100,
				IsPercentage: false,
			}, order.PriceAddon{
				Kind:         "shipping",
				Amount:       60,
				IsPercentage: false,
			}, order.PriceAddon{
				Kind:         "discount",
				Amount:       100,
				IsPercentage: false,
			})

			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[0].ID).Return(&prod, nil)
			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[1].ID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(uint(32), nil)

			mockRepo.On("CreateOrder", &ord2).Return(nil)

			err := s.CreateOrder(&claims, &ord2)
			assert.Nil(t, err)

			assert.Equal(t, float64(910), ord2.Subtotal)
		})

		t.Run("calculate subtotal with items and percentage addons", func(t *testing.T) {
			ord2 := ord
			ord2.Items = []order.Item{
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
			}

			ord2.PriceAddons = append(ord2.PriceAddons, order.PriceAddon{
				Kind:         "fees",
				Amount:       10,
				IsPercentage: true,
			}, order.PriceAddon{
				Kind:         "taxes",
				Amount:       20,
				IsPercentage: true,
			}, order.PriceAddon{
				Kind:         "shipping",
				Amount:       5,
				IsPercentage: true,
			}, order.PriceAddon{
				Kind:         "discount",
				Amount:       30,
				IsPercentage: true,
			})

			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[0].ID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(uint(32), nil)

			mockRepo.On("CreateOrder", &ord2).Return(nil)

			err := s.CreateOrder(&claims, &ord2)
			assert.Nil(t, err)

			assert.Equal(t, float64(306), ord2.Subtotal)
		})

		t.Run("calculate subtotal with items and absolute and percentage addons", func(t *testing.T) {
			ord2 := ord
			ord2.Items = []order.Item{
				{ProductID: prod.ID, VariantID: prod.Variants[0].ID},
			}

			ord2.PriceAddons = append(ord2.PriceAddons, order.PriceAddon{
				Kind:         "fees",
				Amount:       10,
				IsPercentage: true,
			}, order.PriceAddon{
				Kind:         "shipping",
				Amount:       60,
				IsPercentage: false,
			}, order.PriceAddon{
				Kind:         "discount",
				Amount:       50,
				IsPercentage: false,
			})

			mockProdS.On("GetProductByIDAndVariantID", &claims, prod.ID, prod.Variants[0].ID).Return(&prod, nil)
			mockCounterS.On("AddOneToOrder", &claims).Return(uint(32), nil)

			mockRepo.On("CreateOrder", &ord2).Return(nil)

			err := s.CreateOrder(&claims, &ord2)
			assert.Nil(t, err)

			assert.Equal(t, float64(340), ord2.Subtotal)
		})
	})
}
