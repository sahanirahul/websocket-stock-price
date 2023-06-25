package http

import (
	"context"
	"fmt"
	"net/http"
	"sensibull/stocks-api/business/entities/dto"
	"sensibull/stocks-api/business/interfaces/irepo"
	"sensibull/stocks-api/utils/logging"
	"sync"
	"time"
)

type httprepo struct {
}

var once sync.Once
var repo *httprepo

func NewInstrumentHttpRepo() irepo.IInstrumentHttpRepo {
	once.Do(func() {
		repo = &httprepo{}
	})
	return repo
}

func (cr *httprepo) GetUnderLyings(ctx context.Context) ([]dto.Instrument, error) {
	// call 'https://prototype.sbulltech.com/api/underlyings' here
	var response struct {
		Status bool             `json:"success"`
		Data   []dto.Instrument `json:"payload"`
	}
	url := "https://prototype.sbulltech.com/api/underlyings"
	httpreq := HttpRequest{URL: url, Body: nil, Timeout: 2 * time.Second, Method: http.MethodGet}
	status, err := httpreq.InitiateHttpCall(ctx, &response)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_underlyings_http_request", logging.ErrorLevel, logging.Fields{"error": err})
		return nil, err
	}
	if status != http.StatusOK {
		logging.Logger.WriteLogs(ctx, "status_code_not_ok", logging.ErrorLevel, logging.Fields{"statusCode": status})
	}
	return response.Data, nil
}

func (cr *httprepo) GetUnderLyingDerivatives(ctx context.Context, underLyingToken int64) ([]dto.Instrument, error) {
	// call 'https://prototype.sbulltech.com/api/derivatives/{underlying_token}' here
	var response struct {
		Status bool             `json:"success"`
		Data   []dto.Instrument `json:"payload"`
	}
	url := fmt.Sprintf("https://prototype.sbulltech.com/api/derivatives/%d", underLyingToken)
	httpreq := HttpRequest{URL: url, Body: nil, Timeout: 2 * time.Second, Method: http.MethodGet}
	status, err := httpreq.InitiateHttpCall(ctx, &response)
	if err != nil {
		logging.Logger.WriteLogs(ctx, "error_fetching_derivative_http_request", logging.ErrorLevel, logging.Fields{"error": err})
		return nil, err
	}
	if status != http.StatusOK {
		logging.Logger.WriteLogs(ctx, "status_code_not_ok", logging.ErrorLevel, logging.Fields{"statusCode": status})
	}
	return response.Data, nil
}
