package exchange

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func generateResponse(w http.ResponseWriter, base string, target string, rate float64, errorMsg string) {
	_ = json.NewEncoder(w).Encode(RateResponse{
		Base:     base,
		Target:   target,
		Rate:     rate,
		ErrorMsg: errorMsg,
	})
}

func TestGetRateSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/convert" {
			from := r.URL.Query().Get("from")
			to := r.URL.Query().Get("to")

			fromToUSD, fromExist := UsdRates[from]
			toToUSD, toExist := UsdRates[to]
			if !fromExist || !toExist {
				w.WriteHeader(http.StatusBadRequest)
				generateResponse(w, from, to, -1, "invalid currency pair")
				return
			}

			fromInUSD := 1 / fromToUSD
			rate := fromInUSD * toToUSD

			w.WriteHeader(http.StatusOK)
			generateResponse(w, from, to, rate, "")
			return
		}

		w.WriteHeader(http.StatusNotFound)
		generateResponse(w, "", "", -1, "Page not found")
	}))
	defer server.Close()

	exServ := NewExchangeService(server.URL)

	rate, err := exServ.GetRate("USD", "KZT")

	assert.NoError(t, err)
	assert.Greater(t, rate, 0.0)
}

func TestGetRateApiError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/convert" {
			from := r.URL.Query().Get("from")
			to := r.URL.Query().Get("to")

			fromToUSD, fromExist := UsdRates[from]
			toToUSD, toExist := UsdRates[to]
			if !fromExist || !toExist {
				w.WriteHeader(http.StatusBadRequest)
				generateResponse(w, from, to, -1, "invalid currency pair")
				return
			}

			fromInUSD := 1 / fromToUSD
			rate := fromInUSD * toToUSD

			w.WriteHeader(http.StatusOK)
			generateResponse(w, from, to, rate, "")
			return
		}

		w.WriteHeader(http.StatusNotFound)
		generateResponse(w, "", "", -1, "Page not found")
	}))
	defer server.Close()

	exServ := NewExchangeService(server.URL)

	_, err := exServ.GetRate("NOT_CURRENCY", "KZT")

	assert.ErrorContains(t, err, "api error:")
}

func TestGetRateMalformedJson(t *testing.T) {
	wrongJsonServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/convert" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Internal server Error"))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		generateResponse(w, "", "", -1, "Page not found")
	}))
	defer wrongJsonServer.Close()

	exServ := NewExchangeService(wrongJsonServer.URL)

	_, err := exServ.GetRate("USD", "KZT")

	assert.ErrorContains(t, err, "decode error:")
}

func TestGetRateSlowResponse(t *testing.T) {
	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)

		w.WriteHeader(http.StatusOK)

		w.Write([]byte(`{"base":"USD","target":"KZT","rate":469.69}`))
	}))
	defer slowServer.Close()

	exServ := NewExchangeService(slowServer.URL)
	exServ.HTTPClient.Timeout = 50 * time.Millisecond

	_, err := exServ.GetRate("USD", "KZT")

	assert.ErrorContains(t, err, "network error:")
}

func TestGetRateInternalServerError(t *testing.T) {
	wrong500Server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server failure"}`))
	}))
	defer wrong500Server.Close()

	exServ := NewExchangeService(wrong500Server.URL)

	_, err := exServ.GetRate("USD", "KZT")

	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "api error:")
}

func TestGetRateEmptyBody(t *testing.T) {
	emptyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer emptyServer.Close()

	exServ := NewExchangeService(emptyServer.URL)

	_, err := exServ.GetRate("USD", "KZT")

	assert.ErrorContains(t, err, "decode error:")
}
