package resthandlers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/omareloui/odinls/internal/application/core/product"
	"github.com/omareloui/odinls/internal/errs"
	"github.com/omareloui/odinls/web/views"
)

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)
	prods, err := h.app.ProductService.GetCurrentMerchantProducts(claims)
	if err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.ProductsPage(claims, prods))
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) error {
	claims, _ := h.getAuthFromContext(r)

	err := r.ParseForm()
	if err != nil {
		return err
	}

	prod, err := mapFormToProduct(r.PostForm)
	if err != nil {
		return err
	}

	err = h.app.ProductService.CreateProduct(claims, prod)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateProductForm(prod, mapProductToFormData(prod, &valerr), claims.HourlyRate()))
		}

		return err
	}

	if err := renderToBody(w, r, views.ProductOOB(prod, claims.HourlyRate())); err != nil {
		return err
	}
	return respondWithTemplate(w, r, http.StatusOK, views.CreateProductForm(&product.Product{}, &views.ProductFormData{Variants: []views.ProductVariantFormData{{}}}, claims.HourlyRate()))
}

func (h *handler) GetProduct(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)
		prod, err := h.app.ProductService.GetProductByID(claims, id)
		if err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.Product(prod, claims.HourlyRate()))
	}
}

func (h *handler) GetEditProduct(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)
		prod, err := h.app.ProductService.GetProductByID(claims, id)
		if err != nil {
			return err
		}
		return respondWithTemplate(w, r, http.StatusOK, views.EditProduct(prod, mapProductToFormData(prod, &errs.ValidationError{}), claims.HourlyRate()))
	}
}

func (h *handler) EditProduct(id string) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		claims, _ := h.getAuthFromContext(r)

		err := r.ParseForm()
		if err != nil {
			return err
		}

		prod, err := mapFormToProduct(r.PostForm)
		if err != nil {
			return err
		}

		err = h.app.ProductService.UpdateProductByID(claims, id, prod)
		if err != nil {
			if valerr, ok := err.(errs.ValidationError); ok {
				return respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.EditProduct(prod, mapProductToFormData(prod, &valerr), claims.HourlyRate()))
			}
			return err
		}

		return respondWithTemplate(w, r, http.StatusOK, views.Product(prod, claims.HourlyRate()))
	}
}

func mapFormToProduct(form url.Values) (*product.Product, error) {
	prod := &product.Product{
		Name:        form["name"][0],
		Description: form["description"][0],
		Category:    form["category"][0],
		Variants:    []product.Variant{},
	}

	re := regexp.MustCompile(`variant_([\w_]+)-(\d+)`)
	for k, v := range form {
		val := v[0]
		isVariant := strings.Contains(k, "variant_")
		if !isVariant {
			continue
		}

		matches := re.FindStringSubmatch(k)
		key := matches[1]
		idx, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, errs.ErrInvalidNumber
		}

		if len(prod.Variants) < idx+1 {
			for range idx + 1 - len(prod.Variants) {
				prod.Variants = append(prod.Variants, product.Variant{})
			}
		}

		switch key {
		case "id":
			prod.Variants[idx].ID = val
		case "name":
			prod.Variants[idx].Name = val
		case "suffix":
			prod.Variants[idx].Suffix = val
		case "description":
			prod.Variants[idx].Description = val
		case "price":
			if val != "" {
				prod.Variants[idx].Price, err = strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errs.ErrInvalidFloat
				}
			}
		case "wholesale_price":
			if val != "" {
				prod.Variants[idx].WholesalePrice, err = strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errs.ErrInvalidFloat
				}
			}
		case "time_to_craft":
			if val != "" {
				mins, err := strconv.Atoi(val)
				if err != nil {
					return nil, errs.ErrInvalidFloat
				}
				prod.Variants[idx].TimeToCraft = time.Duration(mins * int(time.Minute))
			}
		case "materials_cost":
			if val != "" {
				prod.Variants[idx].MaterialsCost, err = strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errs.ErrInvalidFloat
				}
			}
		}
	}

	return prod, nil
}

func mapProductToFormData(prod *product.Product, valerr *errs.ValidationError) *views.ProductFormData {
	formdata := &views.ProductFormData{
		Name:        views.FormInputData{Value: prod.Name, Error: valerr.Errors.MsgFor("Name")},
		Description: views.FormInputData{Value: prod.Description, Error: valerr.Errors.MsgFor("Description")},
		Category:    views.FormInputData{Value: prod.Category, Error: valerr.Errors.MsgFor("Category")},
		Variants:    []views.ProductVariantFormData{},
	}

	for i, variant := range prod.Variants {
		formdata.Variants = append(formdata.Variants, views.ProductVariantFormData{
			ID:             views.FormInputData{Value: variant.ID, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].ID", i))},
			Name:           views.FormInputData{Value: variant.Name, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Name", i))},
			Description:    views.FormInputData{Value: variant.Description, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Description", i))},
			Suffix:         views.FormInputData{Value: variant.Suffix, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Suffix", i))},
			MaterialsCost:  views.FormInputData{Value: strconv.FormatFloat(variant.MaterialsCost, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].MaterialsCost", i))},
			TimeToCraft:    views.FormInputData{Value: strconv.Itoa(int(variant.TimeToCraft.Minutes())), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].TimeToCraft", i))},
			Price:          views.FormInputData{Value: strconv.FormatFloat(variant.Price, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Price", i))},
			WholesalePrice: views.FormInputData{Value: strconv.FormatFloat(variant.WholesalePrice, 'f', -1, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].WholesalePrice", i))},
		})
	}

	return formdata
}
