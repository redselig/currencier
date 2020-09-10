package controllers

import (
	"context"
	"encoding/xml"
	"github.com/pkg/errors"
	"github.com/redselig/currencier/internal/domain/entity"
	"io"
	"net/http"
	"time"
)

const (
	ErrLoad = "can't pull currency prices from site: %v"
	ErrXML  = "can't extract xml data"
)

var _ entity.CurrencyExternalRepository = (*HTTPClient)(nil)

type HTTPClient struct {
	url    string
	client *http.Client
}

func NewHTTPClient(url string, timeout time.Duration) *HTTPClient {
	client := &http.Client{Timeout: timeout * time.Second}
	return &HTTPClient{
		url:    url,
		client: client}
}
func (hc *HTTPClient) Load(ctx context.Context) ([]*entity.Currency, error) {
	req, err := http.NewRequest(
		"GET", hc.url, nil,
	)
	if err != nil {
		return nil, errors.Wrapf(err, ErrLoad, hc.url)
	}

	req.Header.Add("Accept", "text/html")
	req.Header.Add("User-Agent", "MSIE/15.0")
	req = req.WithContext(ctx)

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, ErrLoad, hc.url)
	}
	defer resp.Body.Close()

	cs, err := XMLExtract(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, ErrLoad, hc.url)
	}
	return cs, nil
}

func XMLExtract(rc io.ReadCloser) ([]*entity.Currency, error) {
	decoded := xml.NewDecoder(rc)
	vals := ValCurs{}
	err := decoded.Decode(&vals)
	if err != nil {
		return nil, errors.Wrapf(err, ErrXML)
	}
	if len(vals.Valute) == 0 {
		return nil, errors.Wrapf(err, ErrXML)
	}
	cs:=XMLValutesToCurrencies(vals.Valute)
	return cs,nil
}

func XMLValutesToCurrencies(vls []Valute) []*entity.Currency {
	cs:=[]*entity.Currency{}
	for _, valute := range vls {
		c:=entity.Currency{
			ID:       valute.ID,
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Nominal:  valute.Nominal,
			Name:     valute.Name,
			Value:    valute.Value,
		}
		cs = append(cs, &c)
	}
	return cs
}

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}
type Valute struct {
	ID       string  `xml:"ID"`
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  int     `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    float64 `xml:"Value"`
}
