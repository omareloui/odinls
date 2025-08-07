package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/a-h/templ"
	jwtadapter "github.com/omareloui/odinls/internal/adapters/jwt"
	"github.com/omareloui/odinls/internal/application/core/client"
	"github.com/omareloui/odinls/internal/application/core/order"
	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetOrders(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims, _ := h.getAuthFromContext(r)
	ords, err := h.app.OrderService.GetOrders(claims)
	if err != nil {
		return err
	}

	prods, clients, err := h.getProdsAndClients(claims)
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.OrdersPage(claims, prods, clients, ords))
}

func (h *handler) CreateOrder(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	claims, err := h.getAuthFromContext(r)
	if err != nil {
		return err
	}

	ord, err := mapFormToOrder(r)
	if err != nil {
		return err
	}

	err = h.app.OrderService.CreateOrder(claims, ord)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			prods, clients, err := h.getProdsAndClients(claims)
			if err != nil {
				return err
			}
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateOrderForm(ord, prods, clients, mapOrderToFormData(ord, &valerr)))
		}
		return err
	}

	prods, clients, err := h.getProdsAndClients(claims)
	if err != nil {
		return err
	}

	err = renderToBody(w, r, views.OrderOOB(ord))
	if err != nil {
		return err
	}

	return respondWithTemplate(w, r, http.StatusOK,
		views.CreateOrderForm(&order.Order{}, prods, clients,
			views.NewDefaultOrderFormData()))
}

func (h *handler) GetOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		claims, err := h.getAuthFromContext(r)
		if err != nil {
			return err
		}

		ord, err := h.app.OrderService.GetOrderByID(claims, id)
		// order.WithPopulatedClient,
		// order.WithPopulatedCraftsman,
		// order.WithPopulatedItemProducts)
		if err != nil {
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.Order(ord))
	}
}

func (h *handler) GetEditOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		claims, err := h.getAuthFromContext(r)
		if err != nil {
			return err
		}

		ord, err := h.app.OrderService.GetOrderByID(claims, id)
		if err != nil {
			return err
		}

		prods, clients, err := h.getProdsAndClients(claims)
		if err != nil {
			return err
		}

		ordFormdata := mapOrderToFormData(ord, &errs.ValidationError{})

		return respondWithTemplate(w, r, http.StatusOK,
			views.EditOrder(ord, prods, clients, ordFormdata))
	}
}

func (h *handler) EditOrder(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
		claims, err := h.getAuthFromContext(r)
		if err != nil {
			return err
		}

		ord, err := mapFormToOrder(r)
		if err != nil {
			return err
		}

		err = h.app.OrderService.UpdateOrderByID(claims, id, ord)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				prods, clients, err := h.getProdsAndClients(claims)
				if err != nil {
					return err
				}
				ordFormdata := mapOrderToFormData(ord, &valerr)
				return respondWithTemplate(w, r, http.StatusUnprocessableEntity,
					views.EditOrder(ord, prods, clients, ordFormdata))
			}
			return err
		}

		return h.GetOrder(id)(w, r)
	}
}

func (h *handler) getProdsAndClients(claims *jwtadapter.AccessClaims) ([]product.Product, []client.Client, error) {
	prods, err := h.app.ProductService.GetProducts(claims)
	if err != nil {
		return nil, nil, err
	}
	clients, err := h.app.ClientService.GetClients(claims)
	if err != nil {
		return nil, nil, err
	}

	return prods, clients, nil
}

func mapFormToOrder(r *http.Request) (*order.Order, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	f := r.PostForm

	o := &order.Order{
		ClientID: f["client_id"][0],
		Status:   f["status"][0],
		Note:     f["note"][0],
		Timeline: order.Timeline{},
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
				items[idx].VariantID = val
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

func mapOrderToFormData(ord *order.Order, valerr *errs.ValidationError) *views.OrderFormData {
	formdata := &views.OrderFormData{
		ClientID: views.FormInputData{Value: ord.ClientID, Error: valerr.Errors.MsgFor("ClientID")},
		Status:   views.FormInputData{Value: ord.Status, Error: valerr.Errors.MsgFor("Status")},
		Timeline: views.OrderTimelineFormData{
			IssuanceDate: views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.IssuanceDate), Error: valerr.Errors.MsgFor("Timeline.IssuanceDate")},
			DueDate:      views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.DueDate), Error: valerr.Errors.MsgFor("Timeline.DueDate")},
			Deadline:     views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.Deadline), Error: valerr.Errors.MsgFor("Timeline.Deadline")},
			DoneOn:       views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.DoneOn), Error: valerr.Errors.MsgFor("Timeline.DoneOn")},
			ShippedOn:    views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.ShippedOn), Error: valerr.Errors.MsgFor("Timeline.ShippedOn")},
			ResolvedOn:   views.FormInputData{Value: formatDateOnlyIfNonZero(ord.Timeline.ResolvedOn), Error: valerr.Errors.MsgFor("Timeline.ResolvedOn")},
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
			Product:     views.FormInputData{Value: item.ProductID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].ProductID", i))},
			Variant:     views.FormInputData{Value: item.VariantID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].VariantID", i))},
			CustomPrice: views.FormInputData{Value: formatFloatIfNonZero(item.CustomPrice), Error: valerr.Errors.MsgFor(fmt.Sprintf("Items[%d].CustomPrice", i))},
			Quantity:    views.FormInputData{Value: "1", Error: ""},
		})
	}

	for i, addon := range ord.PriceAddons {
		formdata.PriceAddons = append(formdata.PriceAddons, views.PriceAddonFormData{
			Kind:         views.FormInputData{Value: addon.Kind, Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].Kind", i))},
			Amount:       views.FormInputData{Value: formatFloatIfNonZero(addon.Amount), Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].Amount", i))},
			IsPercentage: views.FormInputData{Value: formatBooleanIfNonZero(addon.IsPercentage), Error: valerr.Errors.MsgFor(fmt.Sprintf("PriceAddons[%d].IsPercentage", i))},
		})
	}

	for i, recieved := range ord.ReceivedAmounts {
		formdata.ReceivedAmounts = append(formdata.ReceivedAmounts, views.ReceivedAmountFormData{
			Amount: views.FormInputData{Value: formatFloatIfNonZero(recieved.Amount), Error: valerr.Errors.MsgFor(fmt.Sprintf("ReceivedAmounts[%d].Amount", i))},
			Date:   views.FormInputData{Value: formatDateOnlyIfNonZero(recieved.Date), Error: valerr.Errors.MsgFor("ReceivedAmounts[%d].Date")},
		})
	}

	return formdata
}

func setOrderDate(f *url.Values, key string, t *time.Time) (templ.Component, error) {
	var err error
	val := (*f)[key]
	if val != nil {
		*t, err = parseDateOnlyIfExists(val[0])
	}
	return err
}
