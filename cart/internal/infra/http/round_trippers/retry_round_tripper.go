package round_trippers

import "net/http"

type RetryRoundTripper struct {
	rt         http.RoundTripper
	retryLimit uint
}

func NewRetryRoundTripper(rt http.RoundTripper, retryLimit uint) *RetryRoundTripper {
	return &RetryRoundTripper{rt: rt, retryLimit: retryLimit}
}

func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	for retryCount := uint(0); retryCount <= rt.retryLimit; retryCount++ {
		resp, err = rt.rt.RoundTrip(req)
		if err != nil {
			return resp, err
		}

		if resp.StatusCode != 420 && resp.StatusCode != http.StatusTooManyRequests {
			return resp, err
		}
	}

	return resp, err
}
