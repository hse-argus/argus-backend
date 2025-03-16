package client_service

import (
	"argus-backend/internal/logger"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type ClientService struct{}

func NewClientService() *ClientService {
	return &ClientService{}
}

func (cs *ClientService) SendRequest(host string, port int) (int, error) {
	url := fmt.Sprintf("http://%s", host)
	resp, err := resty.New().R().Get(url)
	if err != nil {
		logger.Error(err.Error())
		return -1, err
	}

	return resp.StatusCode(), nil
}
