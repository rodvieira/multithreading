package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type BrasilApiPostalCode struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ViaCepPostalCode struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

const brasil_api_base_url string = "https://brasilapi.com.br/api/cep/v1/"
const viacep_base_url string = "http://viacep.com.br/ws/"
const viacep_format string = "/json"
const desired_postal_code string = "86050070"

func main() {
	c1 := make(chan *BrasilApiPostalCode)
	c2 := make(chan *ViaCepPostalCode)

	go func() {
		brasilApiResp, err := handleGetBrasilApi(desired_postal_code)

		if err != nil {
			fmt.Printf("error in BrasilApi: %d", err)
			return
		}
		c1 <- brasilApiResp

	}()

	go func() {
		viaCepResp, err := handleGetViaCep(desired_postal_code)

		if err != nil {
			fmt.Printf("error in ViaCep: %d", err)
			return
		}
		c2 <- viaCepResp
	}()

	select {
	case cep := <-c1:
		convertedPostalCode, err := json.Marshal(cep)

		if err != nil {
			fmt.Printf("Error to convert BrasilApi JSON: %d", err)
			return
		}

		fmt.Printf("BrasilApi response: %s\n", string(convertedPostalCode))
	case cep := <-c2:
		convertedPostalCode, err := json.Marshal(cep)

		if err != nil {
			fmt.Printf("Error to convert ViaCep JSON: %d", err)
			return
		}

		fmt.Printf("ViaCep response: %s\n", string(convertedPostalCode))
	case <-time.After(time.Second * 1):
		println("Timeout!")
	}
}

func handleGetBrasilApi(postalCode string) (*BrasilApiPostalCode, error) {
	var baseUrl = brasil_api_base_url + postalCode
	resp, err := http.Get(baseUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var postal_code BrasilApiPostalCode
	err = json.Unmarshal(body, &postal_code)
	if err != nil {
		return nil, err
	}
	return &postal_code, nil
}

func handleGetViaCep(postalCode string) (*ViaCepPostalCode, error) {
	var baseUrl = viacep_base_url + postalCode + viacep_format
	resp, err := http.Get(baseUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var postal_code ViaCepPostalCode
	err = json.Unmarshal(body, &postal_code)
	if err != nil {
		return nil, err
	}
	return &postal_code, nil
}
