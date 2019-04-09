package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.alipay-inc.com/ant_agent/metrics-apiserver/pkg/client"
)

func main() {
	http.HandleFunc("/private_api/metric/dataQuery", dataQueryHandler)
	http.ListenAndServe(":8080", nil)
}

func dataQueryHandler(w http.ResponseWriter, r *http.Request) {
	tenantName := r.URL.Query().Get("tenantId")
	workspaceName := r.URL.Query().Get("workspaceId")
	fmt.Printf("tenant name: %s, workspace name: %s\n", tenantName, workspaceName)
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	var payload client.MonitorQueryBody
	fmt.Printf("payload: %s\n", string(b))
	err = json.Unmarshal(b, &payload)
	if err != nil {
		fmt.Println(err)
	}

	respBody := client.APIResponse{
		Data: &client.APIResponseData{
			Data: &client.APIResponseData2{
				IsSuccessed: true,
				Metrics: []*client.Metric{
					{
						Name: "my-metric",
						Labels: map[string]string{
							"key-1": "value-1",
						},
						DataPoints: map[string]float64{
							"1554781416": float64(0.5),
						},
					},
				},
			},
		},
	}
	b, err = json.Marshal(&respBody)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
