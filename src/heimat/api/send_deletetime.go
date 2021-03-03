package api

import (
	"heimatcli/src/x/log"
)

//

// DeleteTime _
// DELETE /api/v1/trackedtimes/203795
func (api *API) DeleteTime(
	timeID int,
) {

	apiURL := api.urlTrackedTimeDelete(timeID)

	resp, _, err := api.httpDelete(api.Token(), apiURL)
	respBody := readBody(resp)
	if err != nil {
		log.Error.Printf(
			"msg='send time delete failed' url='%s'\nerr='%s'\nresp='%s'",
			apiURL,
			err,
			string(respBody),
		)
		return
	}

	if resp.StatusCode >= 300 {
		errorBody := readBody(resp)
		log.Error.Printf("msg='request failed with status'  url='%s'\ncode=%d\nresp=%s", apiURL, resp.StatusCode, errorBody)
		return
	}

}
