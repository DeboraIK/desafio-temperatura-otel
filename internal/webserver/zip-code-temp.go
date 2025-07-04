package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/DeboraIK/lab2-OTEL/internal/dto"
	"github.com/DeboraIK/lab2-OTEL/internal/entity"
	"github.com/DeboraIK/lab2-OTEL/internal/validators"
	usecase "github.com/DeboraIK/lab2-OTEL/use-case"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (s *WebServer) ZipCodeAndTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := s.TemplateData.OTELTracer.Start(ctx, s.TemplateData.RequestNameOTEL+" chamada externa")
	defer span.End()

	var requestData dto.ZipCode
	requestData.ZipCode = r.URL.Query().Get("cep")

	if requestData.ZipCode == "" {
		http.Error(w, entity.ErrZipCodeRequired.Error(), http.StatusUnprocessableEntity)
		return
	}

	if !validators.IsValidZipCode(requestData.ZipCode) {
		http.Error(w, entity.ErrInvalidZipCode.Error(), http.StatusUnprocessableEntity)
		return
	}

	temperatures, err := usecase.Get(ctx, &requestData)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == entity.ErrCannotFindZipcode {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	s.TemplateData.Content = temperatures.City

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(temperatures); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
