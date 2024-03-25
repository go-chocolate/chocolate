package chocomux

import (
	"mime/multipart"
	"net/url"
)

func UnmarshalFormGeneric[T any](form url.Values) (*T, error) {
	var val T
	err := UnmarshalForm(form, &val)
	return &val, err
}

func UnmarshalMultipartFormGeneric[T any](form *multipart.Form) (*T, error) {
	var val T
	err := UnmarshalMultipartForm(form, &val)
	return &val, err
}
