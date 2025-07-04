package webserver

import (
	"encoding/json"
	"net/http"

	"github.com/DeboraIK/lab2-OTEL/internal/dto"
	"github.com/DeboraIK/lab2-OTEL/internal/entity"
	usecase "github.com/DeboraIK/lab2-OTEL/use-case"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func (s *WebServer) ZipCode(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := s.TemplateData.OTELTracer.Start(ctx, s.TemplateData.RequestNameOTEL+" buscando na porta 8080")
	defer span.End()

	if r.Method != http.MethodPost {
		http.Error(w, entity.ErrNotFound.Error(), http.StatusNotFound)
		return
	}

	var requestData dto.ZipCode
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, entity.ErrInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	if requestData.ZipCode == "" {
		http.Error(w, entity.ErrZipCodeRequired.Error(), http.StatusUnprocessableEntity)
		return
	}

	if len(requestData.ZipCode) != 8 {
		http.Error(w, entity.ErrInvalidZipCode.Error(), http.StatusUnprocessableEntity)
		return
	}

	temperature, err := usecase.GetA(ctx, &requestData)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == entity.ErrCannotFindZipcode {
			statusCode = http.StatusNotFound
		}
		http.Error(w, err.Error(), statusCode)
		return
	}

	s.TemplateData.Content = temperature.City

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(temperature); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
