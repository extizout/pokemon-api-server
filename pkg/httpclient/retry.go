package httpclient

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math"
	"net"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

type RetryConfig struct {
	MaxRetries int
	Backoff    time.Duration
}

func DoWithRetry(ctx context.Context, client *http.Client, req *http.Request, cfg RetryConfig) (*http.Response, error) {
	var lastErr error
	var resp *http.Response

	for i := 1; i <= cfg.MaxRetries; i++ {
		clonedReq := req.Clone(ctx)

		resp, lastErr = client.Do(clonedReq)
		log.Infof("HTTP_REQUEST: request to %s", req.URL.String())

		if lastErr == nil {
			bodyString, err := readAndRestoreBody(resp)
			if err != nil {
				log.Warnf("Failed to read response body: %v", err)
			}
			log.Infof("HTTP_RESPONSE: response from %s, STATUS: %d, BODY: %s", req.URL.String(), resp.StatusCode, bodyString)
			return resp, nil
		}

		if !isRetriableError(lastErr) {
			return nil, lastErr
		}

		wait := backoff(cfg.Backoff, i)
		select {
		case <-time.After(wait):
		case <-ctx.Done():
			if errors.Is(lastErr, context.DeadlineExceeded) {
				return nil, ErrRequestTimeout
			}

			netErr := new(net.Error)
			if errors.As(lastErr, netErr) {
				return nil, ErrConnectionError
			}
			return nil, ctx.Err()
		}
	}

	return nil, lastErr
}

func isRetriableError(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr)
}

func backoff(base time.Duration, attempt int) time.Duration {
	//WARN: this is an exponential backoff algorithm
	maxBackoff := 10 * time.Second
	delay := time.Duration(float64(base) * math.Pow(2, float64(attempt)))
	if delay > maxBackoff {
		return maxBackoff
	}
	return delay
}

func readAndRestoreBody(resp *http.Response) (string, error) {
	if resp.Body == nil {
		return "", nil
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes), nil
}
