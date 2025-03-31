package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"

	"route256/cart/internal/models"
)

var getItemTimeout = 5 * time.Second

//go:generate minimock -i=Limiter
type Limiter interface {
	Wait(context.Context) error
}

//nolint:revive
type ProductClient struct {
	client *http.Client
	url    string
	token  string

	limiter Limiter
}

func NewProductClient(client *http.Client, url, token string, limiter Limiter) *ProductClient {
	return &ProductClient{
		client:  client,
		url:     url,
		token:   token,
		limiter: limiter,
	}
}

func (p *ProductClient) GetItem(ctx context.Context, skuID int64) (models.Product, error) {
	url, err := url2.JoinPath(p.url, "product", strconv.FormatInt(skuID, 10))
	if err != nil {
		return models.Product{}, err
	}

	ctx, cancel := context.WithTimeout(ctx, getItemTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return models.Product{}, fmt.Errorf("http new request: %v", err)
	}

	req.Header.Add("X-API-KEY", p.token)

	if err = p.limiter.Wait(ctx); err != nil {
		return models.Product{}, fmt.Errorf("rate limiter: %v", err)
	}
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
