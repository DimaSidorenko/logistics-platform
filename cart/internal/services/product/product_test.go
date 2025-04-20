package product

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"route256/cart/internal/models"
)

//(dosidorenko) goleak не дружит с моей либой трейсинга
//func TestMain(m *testing.M) {
//	goleak.VerifyTestMain(m)
//}

func TestLimiter(t *testing.T) {
	mc := minimock.NewController(t)

	getItemTimeout = 500 * time.Second

	httpClient := &http.Client{
		Transport: &successTransport{},
	}

	limiterMock := NewLimiterMock(mc)

	client := NewProductClient(httpClient, "", "", limiterMock)

	freeTokens := 10
	muWait := sync.Mutex{}
	limiterMock.WaitMock.Set(func(_ context.Context) (err error) {
		muWait.Lock()
		defer muWait.Unlock()

		if freeTokens > 0 {
			freeTokens--
			return nil
		}

		return fmt.Errorf("rate limit exceeded")
	})

	var success atomic.Int32
	eg := errgroup.Group{}
	for i := 0; i < 30; i++ {
		eg.Go(func() error {
			_, err := client.GetItem(context.Background(), 1)
			if nil == err {
				success.Add(1)
			}
			return nil
		})
	}

	_ = eg.Wait()

	require.Equal(t, int32(10), success.Load())
}

type successTransport struct {
	http.RoundTripper
}

func (s *successTransport) RoundTrip(_ *http.Request) (*http.Response, error) {
	temp := models.Product{}

	bodyBytes, err := json.Marshal(temp)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

	resp := &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(bodyBytes)), // Преобразуем байты в io.Reader
		Header:     make(http.Header),
	}

	resp.Header.Set("Content-Type", "application/json")

	return resp, nil
}
