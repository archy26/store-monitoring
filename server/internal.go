package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/store-monitoring/api"
	"github.com/store-monitoring/helper"
	"net/http"
)

func ServeRequest() {
	r := mux.NewRouter()

	// Define the routes
	r.HandleFunc("/trigger_report", api.TriggerReport).Methods("GET")
	r.HandleFunc("/get_report", api.GetReport).Methods("GET")

	// Run the server
	fmt.Println("Starting Server at " + helper.AppConfig.Host + ":" + helper.AppConfig.Port)
	http.ListenAndServe(helper.AppConfig.Host+":"+helper.AppConfig.Port, r)
}
