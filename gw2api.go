package gw2api

import (
	"github.com/go-resty/resty"
	"github.com/pkg/errors"
	"fmt"
)

type GW2v2 struct {
	base string
	key string
}

type apiError struct {
	Text string
}

type Account struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Age          int    `json:"age"`
	World        int    `json:"world"`
	Created      string `json:"created"`
	FractalLevel int    `json:"fractal_level"`
}

type CurrencyAmount struct {
	ID    int `json:"id"`
	Value int `json:"value"`
}

func (v2 *GW2v2) newRequest() *resty.Request {
	return resty.R().
		SetHeader("Authorization", "Bearer " + v2.key).
		SetError(&apiError{})
}

func New(apiKey string) *GW2v2 {
	return &GW2v2{
		base: "https://api.guildwars2.com/v2/",
		key: apiKey,
	}
}

func errMassage(err error, response *resty.Response, method string, endpoint string) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}
	basic := fmt.Sprintf("%s %s", method, endpoint)
	if response == nil {
		return nil, errors.Errorf("%s returned no response?!", basic)
	}
	code := response.StatusCode()
	if code >= 400 {
		if apimsg, ok := response.Error().(*apiError); ok && apimsg != nil {
			return nil, errors.Errorf("%s: a%d: %s", basic, code, apimsg.Text)
		}
		return nil, errors.Errorf("%s: b%d: %s", basic, code, response.String())
	}
	return response, nil
}

func (v2 *GW2v2) Get(endpoint string, result interface{}) (*resty.Response, error) {
	resp, err := v2.newRequest().SetResult(result).Get(v2.base + endpoint)
	if resp, err = errMassage(err, resp, "Get", endpoint); err != nil {
		return nil, err
	}
	return resp, nil
}

func (v2 *GW2v2) GetAccount() (*Account, error) {
	resp, err := v2.Get("account", &Account{})
	if err != nil {
		return nil, err
	}
	if account, ok := resp.Result().(*Account); ok {
		return account, nil
	}
	return nil, errors.Errorf("Could not unmarshal: %v", resp.Result())
}

func (v2 *GW2v2) GetWallet() ([]CurrencyAmount, error) {
	resp, err := v2.Get("account/wallet", make([]CurrencyAmount, 0))
	if err != nil {
		return nil, err
	}
	if wallet, ok := resp.Result().(*[]CurrencyAmount); ok {
		return *wallet, nil
	}
	return nil, errors.Errorf("Could not unmarshal: %v", resp.Result())
}
