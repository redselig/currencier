package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/redselig/currencier/internal/domain/entity"
	"github.com/redselig/currencier/internal/domain/usecase"
	"github.com/redselig/currencier/internal/mocks"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
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
	testLimit  = 3
	testOffset = 1
)

func TestHTTPServer_Serve(t *testing.T) {
	repo := mocks.NewMockRepo(&testCurrency)
	logger := mocks.NewMockLogger()

	currensier:=usecase.NewCurrencierInteractor(nil,repo)
	server:=NewHttpServer("",logger,currensier)
	w := httptest.NewRecorder()
	testCurrencyAnswer,err:=json.Marshal(testCurrency)
	require.Nil(t, err)

	t.Run("GET Currency by id", func(t *testing.T) {
		tCases := []struct {
			title string
			r     *http.Request
			code  int
			body string
		}{
			{"good get",
				httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://fakedns.com/currency/%v", testID), nil),
				200,
				string(testCurrencyAnswer),
			},
			{"bad get",
				httptest.NewRequest(http.MethodGet, "http://fakedns/currency", nil),
				400,
				ErrID,
			},
		}
		for _, tcase := range tCases {
			t.Run(tcase.title, func(t *testing.T) {
				server.getCurrency(w, tcase.r)
				resp := w.Result()
				body,err:=ioutil.ReadAll(resp.Body)
				require.Nil(t, err)
				code:=resp.StatusCode
				require.Equal(t, tcase.code, code)
				require.Equal(t, tcase.body, body)
			})

		}
	})

}
