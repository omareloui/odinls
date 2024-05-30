package resthandlers

import (
	"errors"
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

func (h *handler) GetProducts(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)
	prods, err := h.app.ProductService.GetProducts(claims)
	if err != nil {
		if errors.Is(errs.ErrForbidden, err) {
			respondWithForbidden(w, r)
			return
		}
		respondWithInternalServerError(w, r)
		return
	}
	respondWithTemplate(w, r, http.StatusOK, views.ProductsPage(claims, prods))
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	claims, _ := h.getAuthFromContext(r)

	_ = r.ParseForm()
	prod, err := populateProductFromForm(r.PostForm)
	if err != nil {
		if errors.Is(errs.ErrInvalidNumber, err) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		respondWithInternalServerError(w, r)
		return
	}

	err = h.app.ProductService.CreateProduct(claims, prod)
	if err != nil {
		if valerr, ok := err.(errs.ValidationError); ok {
			respondWithTemplate(w, r, http.StatusUnprocessableEntity, views.CreateProductForm(prod, newProductFormData(prod, &valerr)))
			return
		}

		respondWithInternalServerError(w, r)
		return
	}

	_ = renderToBody(w, r, views.ProductOOB(prod))
	respondWithTemplate(w, r, http.StatusOK, views.CreateProductForm(&product.Product{}, &views.ProductFormData{Variants: []views.ProductVariantFormData{{}}}))
}

func (h *handler) GetProduct(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *handler) GetEditProduct(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (h *handler) EditProduct(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func populateProductFromForm(form url.Values) (*product.Product, error) {
	prod := &product.Product{
		Name:        form["name"][0],
		Description: form["description"][0],
		Category:    form["category"][0],
		Variants:    []product.ProductVariant{},
	}

	re := regexp.MustCompile(`variant_([\w_]+)-([\w_]+)`)
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
				prod.Variants = append(prod.Variants, product.ProductVariant{})
			}
		}

		switch key {
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
					return nil, errs.ErrInvalidNumber
				}
			}
		case "wholesale_price":
			if val != "" {
				prod.Variants[idx].WholesalePrice, err = strconv.ParseFloat(val, 64)
				if err != nil {
					return nil, errs.ErrInvalidNumber
				}
			}
		case "time_to_craft":
			if val != "" {
				mins, err := strconv.Atoi(val)
				if err != nil {
					return nil, errs.ErrInvalidNumber
				}
				prod.Variants[idx].TimeToCraft = time.Duration(mins * int(time.Minute))
			}
		}
	}

	return prod, nil
}

func newProductFormData(prod *product.Product, valerr *errs.ValidationError) *views.ProductFormData {
	formdata := &views.ProductFormData{
		Name:        views.FormInputData{Value: prod.Name, Error: valerr.Errors.MsgFor("Name")},
		Description: views.FormInputData{Value: prod.Description, Error: valerr.Errors.MsgFor("Description")},
		Category:    views.FormInputData{Value: prod.Category, Error: valerr.Errors.MsgFor("Category")},
		Variants:    []views.ProductVariantFormData{},
	}

	for i, variant := range prod.Variants {
		formdata.Variants = append(formdata.Variants, views.ProductVariantFormData{
			Name:           views.FormInputData{Value: variant.Name, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Name", i))},
			Description:    views.FormInputData{Value: variant.Description, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Description", i))},
			Suffix:         views.FormInputData{Value: variant.Suffix, Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Suffix", i))},
			Price:          views.FormInputData{Value: strconv.FormatFloat(variant.Price, 'f', 2, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].Price", i))},
			WholesalePrice: views.FormInputData{Value: strconv.FormatFloat(variant.WholesalePrice, 'f', 2, 64), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].WholesalePrice", i))},
			TimeToCraft:    views.FormInputData{Value: strconv.Itoa(int(variant.TimeToCraft.Minutes())), Error: valerr.Errors.MsgFor(fmt.Sprintf("Variants[%d].TimeToCraft", i))},
		})
	}

	return formdata
}
