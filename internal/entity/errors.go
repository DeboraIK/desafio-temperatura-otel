package entity

import "errors"

var (
	ErrInvalidZipCode        = errors.New("Invalid zipcode")
	ErrCannotFindZipcode     = errors.New("Cannot find zipcode")
	ErrCannotFindCoordinates = errors.New("Cannot find coordinates")
	ErrInternalServer        = errors.New("Internal server error")
	ErrZipCodeRequired       = errors.New("Zipcode required")
	ErrNotFound              = errors.New("Not Found")
	ErrViaCep                = errors.New("Erro com ViaCep")
	ErrGeoAPI                = errors.New("Erro com GeoAPI")
	ErrOpenMeteo             = errors.New("Erro com OpenMeteo")
)
