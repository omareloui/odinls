package resthandlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/errs"
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
	err := r.ParseForm()
	if err != nil {
		return err
	}

	ord, err := mapFormToOrder(r.PostForm)
	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", ord)

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

func mapFormToOrder(f url.Values) (*order.Order, error) {
	var err error

	o := &order.Order{
		ClientID: f["client_id"][0],
		Status:   f["status"][0],
		Note:     f["note"][0],
		Timeline: order.Timeline{},
	}

	o.CustomPrice, err = parseFloatIfExists(f["custom_price"][0])
	if err != nil {
		return nil, err
	}

	o.Timeline.IssuanceDate, err = parseDateOnlyIfExists(f["issuance_date"][0])
	if err != nil {
		return nil, err
	}
	o.Timeline.DueDate, err = parseDateOnlyIfExists(f["due_date"][0])
	if err != nil {
		return nil, err
	}
	o.Timeline.Deadline, err = parseDateOnlyIfExists(f["deadline"][0])
	if err != nil {
		return nil, err
	}
	o.Timeline.DoneOn, err = parseDateOnlyIfExists(f["done_on"][0])
	if err != nil {
		return nil, err
	}
	o.Timeline.ResolvedOn, err = parseDateOnlyIfExists(f["resolved_on"][0])
	if err != nil {
		return nil, err
	}
	o.Timeline.ShippedOn, err = parseDateOnlyIfExists(f["shipped_on"][0])
	if err != nil {
		return nil, err
	}

	multipleKeyFormRegexp := regexp.MustCompile(`^(?:item|addon)_([\w_]+)-(\d+)`)

	type itemsQuantity struct {
		idx      int
		quantity int
	}

	itemsQuantities := []itemsQuantity{}
	items := []order.Item{}

	for k, v := range f {
		val := v[0]
		isItem := strings.Contains(k, "item_")
		isAddon := strings.Contains(k, "addon_")

		if !isItem && !isAddon {
			continue
		}

		matches := multipleKeyFormRegexp.FindStringSubmatch(k)
		key := matches[1]

		idx, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, errs.ErrInvalidFloat
		}

		if isItem {
			if len(items) < idx+1 {
				for range idx + 1 - len(items) {
					items = append(items, order.Item{})
				}
			}

			switch key {
			case "id":
				items[idx].ID = val
			case "product":
				items[idx].ProductID = val
			case "variant":
				items[idx].ProductID = val
			case "custom_price":
				items[idx].CustomPrice, err = parseFloatIfExists(val)
				if err != nil {
					return nil, err
				}
			case "quantity":
				quantity, err := parseIntIfExists(val)
				if err != nil {
					return nil, err
				}
				itemsQuantities = append(itemsQuantities, itemsQuantity{idx: idx, quantity: quantity})
			}
		}

		if isAddon {
			if len(o.PriceAddons) < idx+1 {
				for range idx + 1 - len(o.PriceAddons) {
					o.PriceAddons = append(o.PriceAddons, order.PriceAddon{})
				}
			}

			switch key {
			case "kind":
				o.PriceAddons[idx].Kind = val
			case "amount":
				o.PriceAddons[idx].Amount, err = parseFloatIfExists(val)
				if err != nil {
					return nil, err
				}
			case "is_percentage":
				o.PriceAddons[idx].IsPercentage = true
			}
		}
	}

	for _, iq := range itemsQuantities {
		quantity := iq.quantity
		if quantity < 1 {
			quantity = 1
		}
		for range quantity {
			o.Items = append(o.Items, items[iq.idx])
		}
	}

	return o, nil
}
