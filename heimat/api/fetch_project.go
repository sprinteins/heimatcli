package api

import (
	"encoding/json"
	"heimatcli/heimat"
	"heimatcli/x/log"
	"time"
)

// FetchProjects _
// https://heimat.sprinteins.com/api/v1/employees/29/projects?date=2020-01-15
// {"projects":[{"id":0,"name":"Interne Themen","tasks":[{"id":168,"name":"Becoming Mum"},{"id":76,"name":"Domain Agile"},{"id":151,"name":"Domain Ai"},{"id":72,"name":"Domain Development"},{"id":155,"name":"Domain Ready"},{"id":74,"name":"Domain Vision"},{"id":75,"name":"Domain Zertifizierung"},{"id":21,"name":"Einarbeitung neuer MA"},{"id":103,"name":"Heimat Review"},{"id":243,"name":"ISO Zertifizierung"},{"id":230,"name":"Interne Events"},{"id":257,"name":"Leaders Club DL Edition"},{"id":194,"name":"Mitarbeitergespräch"},{"id":154,"name":"Recruiting Unterstützung"},{"id":31,"name":"Vertrieb"},{"id":38,"name":"Weiterbildung"}]},{"id":34,"name":"ISTDaWo","tasks":[{"id":169,"name":"Entwicklung"},{"id":273,"name":"WISO Entwicklung"}]}]}
func (api *API) FetchProjects() []heimat.Project {

	url := api.urlProjects(api.UserID(), time.Now())
	resp, _, err := httpGet(api.Token(), url, nil)

	if err != nil {
		log.Error.Printf("could not fetch projects: err=%s", err)
	}

	if resp.StatusCode >= 300 {
		log.Error.Printf("could not fetch project, http_status=%d", resp.StatusCode)
	}

	respBodyBytes := readBody(resp)
	projectResponse := &ProjectResponse{}
	err = json.Unmarshal(respBodyBytes, projectResponse)
	if err != nil {
		log.Error.Printf("Could not unmarshal response body: %s\n", err.Error())
	}

	return projectResponse.Projects

}

// ProjectResponse _
type ProjectResponse struct {
	Projects []heimat.Project `json:"projects"`
}
