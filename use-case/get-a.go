package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DeboraIK/lab2-OTEL/internal/dto"
	"github.com/DeboraIK/lab2-OTEL/internal/entity"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func GetA(ctx context.Context, data *dto.ZipCode) (*dto.Temperature, error) {
	url := fmt.Sprintf("http://api-b:8080/temperature?cep=%s", data.ZipCode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("URL: %s\n", url)
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, entity.ErrCannotFindZipcode
	}

	var t dto.Temperature
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrInternalServer
	}

	return &t, nil
}
