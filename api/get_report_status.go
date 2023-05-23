package api

import (
	"fmt"
	"github.com/store-monitoring/db"
	"github.com/store-monitoring/helper"
	"io"
	"net/http"
	"os"
)

type ReportStatus string

const (
	ReportStatusRunning  ReportStatus = "Running"
	ReportStatusComplete ReportStatus = "Completed"
)

func GetReport(w http.ResponseWriter, r *http.Request) {
	reportID := r.URL.Query().Get("report_id")
	reportStatus := db.ReportStatus{}
	helper.DB.Where("report_id = ?", reportID).Find(&reportStatus)
	fmt.Println(reportStatus, reportID)
	if ReportStatus(reportStatus.Status) == ReportStatusComplete {
		csvFile, err := os.Open("data/" + reportID + ".csv")
		if err != nil {
			fmt.Println(err)
			http.NotFound(w, r)
		}
		defer csvFile.Close()

		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename="+reportID+".csv")

		io.Copy(w, csvFile)
	} else if ReportStatus(reportStatus.Status) == ReportStatusRunning {
		fmt.Fprintf(w, "Running")
	} else {
		http.NotFound(w, r)
	}
}
