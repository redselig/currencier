package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"

	"github.com/redselig/currencier/internal/domain/entity"
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
	decoded.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	vals := ValCurs{}
	err := decoded.Decode(&vals)
	if err != nil {
		return nil, errors.Wrapf(err, ErrXML)
	}
	if len(vals.Valute) == 0 {
		return nil, errors.Wrapf(err, ErrXML)
	}
	cs, err := XMLValutesToCurrencies(vals.Valute)
	if err != nil {
		return nil, errors.Wrapf(err, ErrXML)
	}
	return cs, nil
}

func XMLValutesToCurrencies(vls []Valute) ([]*entity.Currency, error) {
	cs := []*entity.Currency{}
	for _, valute := range vls {
		rate, err := strconv.ParseFloat(strings.Replace(valute.Value, ",", ".", 1), 64)
		if err != nil {
			return nil, err
		}
		rate = math.Round(rate*100) / 100
		c := entity.Currency{
			ID:       valute.ID,
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
			Nominal:  valute.Nominal,
			Name:     valute.Name,
			Value:    rate,
		}
		cs = append(cs, &c)
	}
	return cs, nil
}

type ValCurs struct {
	Valute []Valute `xml:"Valute"`
}
type Valute struct {
	ID       string `xml:"ID,attr"`
	NumCode  int    `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nominal  int    `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
