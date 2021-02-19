package api

import (
	"encoding/json"
	"heimatcli/src/heimat"
	"heimatcli/src/x/log"
	"time"
)

// SendCreateTime _
// POST https://heimat.sprinteins.com/api/v1/trackedtimes
// {"employee":{"id":29},"date":"2021-02-18","trackedTimes":[{"start":"08:00","end":"10:00","note":"Development","task":{"id":300}}]}
func (api *API) SendCreateTime(
	userID int,
	date time.Time,
	start string,
	end string,
	note string,
	task heimat.Task,
) {

	apiURL := api.urlTrackedTimeCreate()

	createTimePayload := createTrackedTimePayload{
		Date:     NewHeimatDate(date),
		Employee: employeeIDPayload{ID: userID},
		TrackedTimes: []trackedTimePayload{
			{
				Start: start,
				End:   end,
				Note:  note,
				Task:  taskIDPayload{ID: task.ID},
			},
		},
	}

	payload, err := json.Marshal(createTimePayload)
	if err != nil {
		log.Error.Printf("send create time could not marshal payload:%s", err)
		return
	}

	resp, _, err := api.httpPost(api.Token(), apiURL, nil, payload)
	if err != nil {
		log.Error.Printf("send time crate post failed:%s", err)
		return
	}

	if resp.StatusCode >= 300 {
		errorBody := readBody(resp)
		log.Error.Printf("msg='request failed with status' code=%d resp=%s", resp.StatusCode, errorBody)
		return
	}

}

type createTrackedTimePayload struct {
	Date         HeimatDate           `json:"date"`
	Employee     employeeIDPayload    `json:"employee"`
	TrackedTimes []trackedTimePayload `json:"trackedTimes"`
}

type employeeIDPayload struct {
	ID int `json:"id"`
}

type trackedTimePayload struct {
	End   string        `json:"end"`
	Start string        `json:"start"`
	Note  string        `json:"note"`
	Task  taskIDPayload `json:"task"`
}

type taskIDPayload struct {
	ID int `json:"id"`
}
