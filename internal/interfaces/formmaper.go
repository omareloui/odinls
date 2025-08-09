package interfaces

type FormMapper interface {
	MapToForm(doc any, err error, formData any) error
}
