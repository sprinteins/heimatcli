package api

import "fmt"

const (
	urlAuthentication = "/authentication"
	urlEmployees      = "/employees"
)

func (api API) urlAuthentication() string {

	apiURL := urlAuthentication

	fullURL := fmt.Sprintf("%s/%s", api.baseURL.String(), apiURL)
	return fullURL
}

func (api API) urlEmployeeByID(id int) string {

	apiURL := fmt.Sprintf("%s/%d", urlEmployees, id)

	fullURL := fmt.Sprintf("%s/%s", api.baseURL.String(), apiURL)
	return fullURL
}
