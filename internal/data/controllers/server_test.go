package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"

	"github.com/redselig/currencier/internal/domain/entity"
	"github.com/redselig/currencier/internal/domain/usecase"
	"github.com/redselig/currencier/internal/mocks"
)

var (
	testID       = "R01020A"
	testName     = "Азербайджанский манат"
	testRate     = 44.7113
	testCurrency = entity.Currency{
		ID:       testID,
		NumCode:  0,
		CharCode: "",
		Nominal:  1,
		Name:     testName,
		Value:    testRate,
	}
)

func TestHTTPServer_Serve(t *testing.T) {
	repo := mocks.NewMockRepo(&testCurrency)
	logger := mocks.NewMockLogger()

	currensier := usecase.NewCurrencierInteractor(nil, repo)
	server := NewHttpServer("", logger, currensier)
	testCurrencyAnswer, err := json.Marshal(testCurrency)
	require.Nil(t, err)
	testCurrenciesAnswer, err := json.Marshal([]*entity.Currency{&testCurrency})
	require.Nil(t, err)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	t.Run("GET currency by id", func(t *testing.T) {
		tCases := []struct {
			title string
			vars  map[string]string
			code  int
			body  string
		}{
			{"good get",
				map[string]string{"id": testID},
				200,
				string(testCurrencyAnswer),
			},
			{"bad get",
				map[string]string{},
				400,
				ErrID + "\n",
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				w := httptest.NewRecorder()
				req = mux.SetURLVars(req, tcase.vars)

				server.getCurrency(w, req)
				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				require.Nil(t, err)
				code := resp.StatusCode

				require.Equal(t, tcase.code, code)
				require.Equal(t, tcase.body, string(body))
			})

		}
	})
	t.Run("GET all currencies", func(t *testing.T) {
		tCases := []struct {
			title string
			req   *http.Request
			code  int
			body  string
		}{
			{"good get Currencies",
				httptest.NewRequest(http.MethodGet, "/currencies?limit=10&offset=5", nil),
				200,
				string(testCurrenciesAnswer),
			},
			{"bad get Currencies",
				httptest.NewRequest(http.MethodGet, "/currencies?limit=bad_value&offset=5", nil),
				400,
				"strconv.Atoi: parsing \"bad_value\": invalid syntax\n",
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				w := httptest.NewRecorder()

				server.getCurrencies(w, tcase.req)
				resp := w.Result()
				body, err := ioutil.ReadAll(resp.Body)
				require.Nil(t, err)
				code := resp.StatusCode

				require.Equal(t, tcase.code, code)
				require.Equal(t, tcase.body, string(body))
			})

		}
	})
}
