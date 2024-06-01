package resthandlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	ords, err := h.app.OrderService.GetCurrentMerchantOrders(claims)
	if err != nil {
		return err
	}
	prods, clients, err := h.getMerchantProdsAndClients(claims)
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims, prods, clients, ords))
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	claims, err := h.getAuthFromContext(r)
	if err != nil {
		return err
	}

	err = r.ParseForm()
	if err != nil {
		return err
	}

	ord, err := mapFormToOrder(r.PostForm)
	if err != nil {
		return err
	}

	err = h.app.OrderService.CreateOrder(claims, ord)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			prods, clients, err := h.getMerchantProdsAndClients(claims)
			if err != nil {
				return err
			}
			fmt.Println("Validation error found:", valerr.Errors)
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateOrderForm(ord, prods, clients, mapOrderToFormData(ord, &valerr)))
		}
		return err
	}

	fmt.Printf("=============> %+v\n", ord)

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

func (h *handler) getMerchantProdsAndClients(claims *jwtadapter.JwtAccessClaims) ([]product.Product, []client.Client, error) {
	prods, err := h.app.ProductService.GetCurrentMerchantProducts(claims)
	if err != nil {
		return nil, nil, err
	}
	clients, err := h.app.ClientService.GetCurrentMerchantClients(claims)
	if err != nil {
		return nil, nil, err
	}

	return prods, clients, nil
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

	err = setOrderDate(&f, "issuance_date", &o.Timeline.IssuanceDate)
	if err != nil {
		return nil, err
	}
	err = setOrderDate(&f, "due_date", &o.Timeline.DueDate)
	if err != nil {
		return nil, err
	}
	err = setOrderDate(&f, "deadline", &o.Timeline.Deadline)
	if err != nil {
		return nil, err
	}
	err = setOrderDate(&f, "done_on", &o.Timeline.DoneOn)
	if err != nil {
		return nil, err
	}
	err = setOrderDate(&f, "resolved_on", &o.Timeline.ResolvedOn)
	if err != nil {
		return nil, err
	}
	err = setOrderDate(&f, "shipped_on", &o.Timeline.ShippedOn)
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
		isItem := strings.HasPrefix(k, "item_")
		isAddon := strings.HasPrefix(k, "addon_")

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

func setOrderDate(f *url.Values, key string, t *time.Time) error {
	var err error
	val := (*f)[key]
	if val != nil {
		*t, err = parseDateOnlyIfExists(val[0])
	}
	return err
}

func mapOrderToFormData(ord *order.Order, valerr *errs.ValidationError) *views.OrderFormData {
	formdata := &views.OrderFormData{
		// Timeline    TimelineFormData `json:"timeline"`
		// Note        FormInputData    `json:"note"`
		// CustomPrice FormInputData    `json:"custom_price"`

		// Items           []OrderItemFormData      `json:"items"`
		// PriceAddons     []PriceAddonFormData     `json:"price_addons"`
		// ReceivedAmounts []ReceivedAmountFormData `json:"received_amounts"`

		ClientID: views.FormInputData{Value: ord.ClientID, Error: valerr.Errors.MsgFor("ClientID")},
		Status:   views.FormInputData{Value: ord.Status, Error: valerr.Errors.MsgFor("Status")},
		Timeline: views.OrderTimelineFormData{
			IssuanceDate: views.FormInputData{Value: ord.Timeline.IssuanceDate.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.IssuanceDate")},
			DueDate:      views.FormInputData{Value: ord.Timeline.DueDate.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.DueDate")},
			Deadline:     views.FormInputData{Value: ord.Timeline.Deadline.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.Deadline")},
			DoneOn:       views.FormInputData{Value: ord.Timeline.DoneOn.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.DoneOn")},
			ShippedOn:    views.FormInputData{Value: ord.Timeline.ShippedOn.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.ShippedOn")},
			ResolvedOn:   views.FormInputData{Value: ord.Timeline.ResolvedOn.Format(time.DateOnly), Error: valerr.Errors.MsgFor("Timeline.ResolvedOn")},
		},
		Items:           []views.OrderItemFormData{},
		PriceAddons:     []views.PriceAddonFormData{},
		ReceivedAmounts: []views.ReceivedAmountFormData{},
	}

	for i, item := range ord.Items {
		existsIdx := slices.IndexFunc(formdata.Items, func(existsItem views.OrderItemFormData) bool {
			return existsItem.Product.Value == item.ProductID &&
				existsItem.Variant.Value == item.VariantID &&
				existsItem.CustomPrice.Value == strconv.FormatFloat(item.CustomPrice, 'f', -1, 64)
		})
		if existsIdx > -1 {
			currQuantityStr := formdata.Items[existsIdx].Quantity.Value
			currQuantity, _ := strconv.Atoi(currQuantityStr)
			formdata.Items[existsIdx].Quantity.Value = strconv.Itoa(currQuantity + 1)
			continue
		}

		formdata.Items = append(formdata.Items, views.OrderItemFormData{
			ID:          views.FormInputData{Value: item.ID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].ID", i))},
			Product:     views.FormInputData{Value: item.ProductID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].Product", i))},
			Variant:     views.FormInputData{Value: item.VariantID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].Variant", i))},
			CustomPrice: views.FormInputData{Value: strconv.FormatFloat(item.CustomPrice, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].CustomPrice", i))},
			Quantity:    views.FormInputData{Value: "1", Error: ""},
		})
	}

	for i, addon := range ord.PriceAddons {
		// FIXME: this doesn't work with alpin's chebox
		isPercentage := ""
		if addon.IsPercentage {
			isPercentage = "on"
		}
		formdata.PriceAddons = append(formdata.PriceAddons, views.PriceAddonFormData{
			Kind:         views.FormInputData{Value: addon.Kind, Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].Kind", i))},
			Amount:       views.FormInputData{Value: strconv.FormatFloat(addon.Amount, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].Amount", i))},
			IsPercentage: views.FormInputData{Value: isPercentage, Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].IsPercentage", i))},
		})
	}

	for i, recieved := range ord.ReceivedAmounts {
		formdata.ReceivedAmounts = append(formdata.ReceivedAmounts, views.ReceivedAmountFormData{
			Amount: views.FormInputData{Value: strconv.FormatFloat(recieved.Amount, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("ReceivedAmounts[%d].Amount", i))},
			Date:   views.FormInputData{Value: recieved.Date.Format(time.DateOnly), Error: valerr.Errors.MsgFor("ReceivedAmounts[%d].Date")},
		})
	}

	return formdata
}
