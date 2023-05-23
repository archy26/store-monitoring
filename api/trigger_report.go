package api

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/store-monitoring/db"
	"github.com/store-monitoring/helper"
	"github.com/store-monitoring/logic"
	"net/http"
)

var runningReportsChan = make(chan struct{}, 1)

func TriggerReport(w http.ResponseWriter, r *http.Request) {
	reportID := uuid.New()
	select {
	case runningReportsChan <- struct{}{}:
		go func() {
			defer func() {
				if r := recover(); r != nil {
					reportStatus := db.ReportStatus{ReportID: reportID.String(), Status: "Error"}
					helper.DB.Save(&reportStatus)
				}
				<-runningReportsChan
			}()

			reportStatus := db.ReportStatus{ReportID: reportID.String(), Status: "Running"}
			helper.DB.Create(&reportStatus)
			err := logic.BuildReport(reportID.String())
			if err != nil {
				panic(err)
			}
		}()

		fmt.Fprintf(w, "%s\n", reportID.String())

	default:
		fmt.Fprintf(w, "Maximum concurrent reports reached\n")
	}

}
