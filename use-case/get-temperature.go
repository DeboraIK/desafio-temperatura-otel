package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/DeboraIK/lab2-OTEL/internal/dto"
	"github.com/DeboraIK/lab2-OTEL/internal/entity"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type Address struct {
	Code  string `json:"cep"`
	State string `json:"estado"`
	City  string `json:"localidade"`
}

func Get(ctx context.Context, data *dto.ZipCode) (*dto.Temperature, error) {
	address, err := getCityFromZipcode(ctx, data.ZipCode)
	if err != nil {
		return nil, err
	}

	temperature, err := getTemperature(ctx, address.City)
	if err != nil {
		return nil, err
	}

	return temperature, nil
}

func getCityFromZipcode(ctx context.Context, zipcode string) (*Address, error) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", zipcode)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrViaCep
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrViaCep
	}
	defer resp.Body.Close()

	var a Address
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrViaCep
	}

	if a.City == "" {
		return nil, entity.ErrCannotFindZipcode
	}

	return &a, nil
}

func getTemperature(ctx context.Context, city string) (*dto.Temperature, error) {
	geoAPIURL := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=pt&format=json", url.QueryEscape(city))
	respGeo, err := http.Get(geoAPIURL)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrCannotFindCoordinates
	}
	defer respGeo.Body.Close()

	var coordinates struct {
		Results []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"results"`
	}

	bodyGeo, err := io.ReadAll(respGeo.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrGeoAPI
	}

	if err := json.Unmarshal(bodyGeo, &coordinates); err != nil {
		log.Println(err.Error())
		return nil, entity.ErrGeoAPI
	}

	if len(coordinates.Results) == 0 {
		return nil, fmt.Errorf("não foi possível encontrar coordenadas para a cidade: %s", city)
	}

	latitude := coordinates.Results[0].Latitude
	longitude := coordinates.Results[0].Longitude

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&current_weather=true", latitude, longitude)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrOpenMeteo
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrOpenMeteo
	}
	defer resp.Body.Close()

	var temp dto.OpenMeteoResponse
	var d struct {
		Current dto.Temperature `json:"current"`
	}

	err = json.NewDecoder(resp.Body).Decode(&temp)
	if err != nil {
		log.Println(err.Error())
		return nil, entity.ErrOpenMeteo
	}

	d.Current.Temp_C = temp.CurrentWeather.Temperature
	d.Current.Temp_K = d.Current.Temp_C + 273.15
	d.Current.Temp_F = (d.Current.Temp_C * 1.8) + 32
	d.Current.City = city

	return &d.Current, nil
}
