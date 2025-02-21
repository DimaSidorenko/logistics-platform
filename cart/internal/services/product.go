package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	dto2 "route256/cart/internal/usecases/cart/dto"
)

type ProductClient struct {
	client *http.Client
	host   string
	port   int
	url    string
	token  string
}

func NewProductClient(client *http.Client, host string, port int, token string) *ProductClient {
	return &ProductClient{
		client: client,
		host:   host,
		port:   port,
		url:    fmt.Sprintf("http://%s:%d", host, port),
		token:  token,
	}
}

func (p *ProductClient) GetItem(skuID dto2.SkuID) (dto2.Product, error) {
	url := fmt.Sprintf("%s/product/%d", p.url, skuID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dto2.Product{}, fmt.Errorf("http new request: %v", err)
	}

	req.Header.Add("X-API-KEY", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return dto2.Product{}, fmt.Errorf("http do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return dto2.Product{}, dto2.ErrItemNotFound
		}

		return dto2.Product{}, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return dto2.Product{}, fmt.Errorf("http read all: %v", err)
	}

	var product dto2.Product
	if err := json.Unmarshal(body, &product); err != nil {
		return dto2.Product{}, fmt.Errorf("json unmarshal: %v", err)
	}

	return product, nil
}
