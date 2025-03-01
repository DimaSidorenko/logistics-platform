package product

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"route256/cart/internal/models"
	"strconv"
)

//nolint:revive
type ProductClient struct {
	client *http.Client
	url    string
	token  string
}

func NewProductClient(client *http.Client, url, token string) *ProductClient {
	return &ProductClient{
		client: client,
		url:    url,
		token:  token,
	}
}

func (p *ProductClient) GetItem(skuID int64) (models.Product, error) {
	url, err := url2.JoinPath(p.url, "product", strconv.FormatInt(skuID, 10))
	if err != nil {
		return models.Product{}, err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.Product{}, fmt.Errorf("http new request: %v", err)
	}

	req.Header.Add("X-API-KEY", p.token)

	resp, err := p.client.Do(req)
	if err != nil {
		return models.Product{}, fmt.Errorf("http do request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return models.Product{}, models.ErrItemNotFound
		}

		return models.Product{}, errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.Product{}, fmt.Errorf("http read all: %v", err)
	}

	var product models.Product
	if err := json.Unmarshal(body, &product); err != nil {
		return models.Product{}, fmt.Errorf("json unmarshal: %v", err)
	}

	return product, nil
}
