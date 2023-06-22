package http

import (
	"context"
	"net/http"
	"sensibull/stocks-api/business/entities/dto"
	"sensibull/stocks-api/utils/logging"
	"time"
)

type httprepo struct {
}

func (cr *httprepo) GetUnderLyings(ctx context.Context) ([]dto.UnderLyings, error) {
	// call 'https://prototype.sbulltech.com/api/underlyings' here
	var response struct {
		Status bool              `json:"success"`
		Data   []dto.UnderLyings `json:"payload"`
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
