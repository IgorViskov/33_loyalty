package services

import (
	"encoding/json"
	"fmt"
	"github.com/IgorViskov/33_loyalty/internal/api/models"
	"github.com/IgorViskov/33_loyalty/internal/config"
	"github.com/IgorViskov/33_loyalty/internal/domain"
	"github.com/IgorViskov/33_loyalty/internal/domain/statuses"
	"io"
	"net/http"
	"net/url"
	"slices"
)

var invalidStatuses []int = []int{http.StatusNoContent, http.StatusInternalServerError}

type ExternalAccrualService interface {
	GetAccrual(order string) (domain.Accrual, error)
}

type externalAccrualService struct {
	client *http.Client
	url    *url.URL
}

func NewExternalAccrualService(conf *config.AppConfig) ExternalAccrualService {
	return &externalAccrualService{
		client: http.DefaultClient,
		url:    conf.AccrualHost,
	}
}

func (e *externalAccrualService) GetAccrual(order string) (domain.Accrual, error) {
	response, err := e.client.Do(e.request(order))
	if err != nil {
		return domain.Accrual{}, err
	}
	if slices.Contains(invalidStatuses, response.StatusCode) {
		return invalid(order), nil
	}
	jsons, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.Accrual{}, err
	}
	var model models.AccrualResponse
	err = json.Unmarshal(jsons, &model)
	if err != nil {
		return domain.Accrual{}, err
	}
	_ = response.Body.Close()
	return domain.Accrual{
		OrderNumber: order,
		Status:      model.Status,
		Value:       model.Accrual,
	}, nil
}

func (e *externalAccrualService) request(order string) *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/orders/%s", e.url, order), nil)
	if err != nil {
		panic(err)
	}
	return req
}

func invalid(order string) domain.Accrual {
	return domain.Accrual{
		OrderNumber: order,
		Status:      statuses.INVALID,
	}
}
