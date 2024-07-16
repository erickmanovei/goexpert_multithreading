package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ApiBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ApiViaCep struct {
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

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type ApiResponse struct {
	Address Address
	Source  string
	Error   error
}

func fetchFromBrasilAPI(cep string, ch chan<- ApiResponse) {
	url := "https://brasilapi.com.br/api/cep/v1/" + cep
	resp, err := http.Get(url)
	if err != nil {
		ch <- ApiResponse{Error: err}
		return
	}
	defer resp.Body.Close()

	var addressBrasilApi ApiBrasilApi
	err = json.NewDecoder(resp.Body).Decode(&addressBrasilApi)
	if err != nil {
		ch <- ApiResponse{Error: err}
		return
	}

	address := Address{
		Cep:         addressBrasilApi.Cep,
		Logradouro:  addressBrasilApi.Street,
		Complemento: "",
		Bairro:      addressBrasilApi.Neighborhood,
		Localidade:  addressBrasilApi.City,
		Uf:          addressBrasilApi.State,
	}

	ch <- ApiResponse{Address: address, Source: "BrasilAPI"}
}

func fetchFromViaCEP(cep string, ch chan<- ApiResponse) {
	url := "http://viacep.com.br/ws/" + cep + "/json/"
	resp, err := http.Get(url)
	if err != nil {
		ch <- ApiResponse{Error: err}
		return
	}
	defer resp.Body.Close()

	var addressViaCep ApiViaCep
	err = json.NewDecoder(resp.Body).Decode(&addressViaCep)
	if err != nil {
		ch <- ApiResponse{Error: err}
		return
	}

	address := Address{
		Cep:         addressViaCep.Cep,
		Logradouro:  addressViaCep.Logradouro,
		Complemento: addressViaCep.Complemento,
		Bairro:      addressViaCep.Bairro,
		Localidade:  addressViaCep.Localidade,
		Uf:          addressViaCep.Uf,
	}

	ch <- ApiResponse{Address: address, Source: "ViaCEP"}
}

func main() {
	cep := "41650000"
	ch1 := make(chan ApiResponse)
	ch2 := make(chan ApiResponse)

	go fetchFromBrasilAPI(cep, ch1)
	go fetchFromViaCEP(cep, ch2)

	select {
	case res := <-ch1:
		if res.Error != nil {
			fmt.Println("Erro:", res.Error)
		} else {
			fmt.Printf("Endereço: %+v\n", res.Address)
			fmt.Println("API de Origem:", res.Source)
		}
	case res := <-ch2:
		if res.Error != nil {
			fmt.Println("Erro:", res.Error)
		} else {
			fmt.Printf("Endereço: %+v\n", res.Address)
			fmt.Println("API de Origem:", res.Source)
		}
	case <-time.After(1 * time.Second):
		fmt.Println("Erro: Timeout")
	}
}
