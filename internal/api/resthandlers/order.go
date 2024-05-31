package resthandlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	ords, err := h.app.OrderService.GetOrders(claims)
	if err != nil {
		return err
	}
	prods, err := h.app.ProductService.GetCurrentMerchantProducts(claims)
	if err != nil {
		return err
	}
	clients, err := h.app.ClientService.GetCurrentMerchantClients(claims)
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims, prods, clients, ords))
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	ord, err := mapFormToOrder(r)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", ord)
	fmt.Printf("client %s\n", ord.ClientID)
	fmt.Printf("timeline %+v\n", ord.Timeline)

	return errors.New("just fail")
}

func (h *handler) GetOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func (h *handler) GetEditOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func (h *handler) EditOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		return nil
	}
}

func mapFormToOrder(r *http.Request) (*order.Order, error) {
	var err error

	o := &order.Order{
		ClientID: r.FormValue("client_id"),
		Status:   r.FormValue("status"),
		Note:     r.FormValue("note"),
		Timeline: order.Timeline{},
	}

	o.CustomPrice, err = parseFloatIfExists(r.FormValue("custom_price"))
	if err != nil {
		return nil, err
	}

	o.Timeline.IssuanceDate, err = parseDateOnlyIfExists(r.FormValue("issuance_date"))
	if err != nil {
		return nil, err
	}

	return o, nil
}
