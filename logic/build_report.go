package logic

import (
	"encoding/csv"
	"fmt"
	"github.com/store-monitoring/db"
	"github.com/store-monitoring/helper"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

var maxCurrTime time.Time

func BuildReport(reportId string) error {

	result := helper.DB.Model(&db.StoreData{}).Select("MAX(timestamp_utc)").Scan(&maxCurrTime)
	if result.Error != nil {
		return result.Error
	}
	stores := []db.StoreData{}

	fetchAllStores := helper.DB.Find(&stores)
	if fetchAllStores.Error != nil {
		return fetchAllStores.Error
	}
	storesAdded := map[string]bool{}
	rows := [][]string{}
	for _, store := range stores {
		if storesAdded[store.StoreID] {
			continue
		}
		newRow := []string{}
		newRow = append(newRow, store.StoreID)
		timezone := db.Timezone{}
		helper.DB.Where("store_id = ?", store.StoreID).Find(&timezone)
		if timezone.TimezoneStr == "" {
			timezone.TimezoneStr = "America/Chicago"
		}
		timezoneStr, timezoneStrErr := time.LoadLocation(timezone.TimezoneStr)
		if timezoneStrErr != nil {
			fmt.Println("Error loading timezone:", timezoneStrErr)
			return timezoneStrErr
		}

		// Convert UTC time to respective timezone
		newTimezone := maxCurrTime.In(timezoneStr)
		weekday := int(newTimezone.Weekday()) - 1
		if weekday < 0 {
			weekday = 6 // Sunday is converted to 6
		}
		storeStatus := []db.StoreData{}

		helper.DB.Where("store_id = ?", store.StoreID).Find(&storeStatus)
		businessHours := []db.StoreBusinessHours{}
		helper.DB.Where("store_id = ? AND day = ?", store.StoreID, weekday).Find(&businessHours)
		uptime_lasthour, downtime_lasthour := getupdateoftoday(store.StoreID, weekday, businessHours, storeStatus, timezoneStr, maxCurrTime)
		newRow = append(newRow, fmt.Sprintf("%f", (uptime_lasthour/24.0)*60.0))
		newRow = append(newRow, fmt.Sprintf("%f", (downtime_lasthour/24.0)*60.0))
		uptime_lastday, downtime_lastday := getupdateoflastday(store.StoreID, timezoneStr, maxCurrTime)
		newRow = append(newRow, fmt.Sprintf("%f", uptime_lastday))
		newRow = append(newRow, fmt.Sprintf("%f", downtime_lastday))
		uptime_lastweek, downtime_lastweek := getupdateoflastweek(store.StoreID, timezoneStr, maxCurrTime)
		newRow = append(newRow, fmt.Sprintf("%f", uptime_lastweek))
		newRow = append(newRow, fmt.Sprintf("%f", downtime_lastweek))
		rows = append(rows, newRow)
		storesAdded[store.StoreID] = true

	}

	//pushing into file
	err := createFileAndUpdateReportStatus(reportId, rows)
	return err
}

func createFileAndUpdateReportStatus(reportId string, rows [][]string) error {
	file, err := os.Create("data/" + reportId + ".csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	headers := []string{
		"store_id",
		"uptime_last_hour(in minutes)",
		"downtime_last_hour(in minutes)",
		"uptime_last_day(in hours)",
		"downtime_last_day(in hours)",
		"uptime_last_week(in hours)",
		"downtime_last_week(in hours)",
	}

	writer := csv.NewWriter(file)

	// Write headers
	writer.Write(headers)

	// Write rows
	for _, row := range rows {
		writer.Write(row)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil
	}
	reportStatus := db.ReportStatus{ReportID: reportId, Status: "Completed"}
	helper.DB.Save(&reportStatus)
	return nil
}

func getupdateoftoday(storeID string, weekday int, businessHours []db.StoreBusinessHours, storeStatus []db.StoreData, timezone *time.Location, maxCurrTime time.Time) (float64, float64) {
	uptime := 0.0
	downtime := 0.0
	totalTime := 0.0
	sort.Slice(storeStatus, func(i, j int) bool {
		return storeStatus[i].TimestampUTC.Before(storeStatus[j].TimestampUTC)
	})
	if len(businessHours) == 0 {
		businessHours = append(businessHours, db.StoreBusinessHours{StoreID: storeID, Day: weekday, StartTimeLocal: "00:00:00", EndTimeLocal: "23:59:59"})
	}
	for _, businessHour := range businessHours {
		parts := strings.Split(businessHour.StartTimeLocal, ":")
		hour, _ := strconv.Atoi(parts[0])
		minute, _ := strconv.Atoi(parts[1])
		second, _ := strconv.Atoi(parts[2])

		parts1 := strings.Split(businessHour.EndTimeLocal, ":")
		hour1, _ := strconv.Atoi(parts1[0])
		minute1, _ := strconv.Atoi(parts1[1])
		second1, _ := strconv.Atoi(parts1[2])

		openingHour := time.Date(maxCurrTime.Year(), maxCurrTime.Month(), maxCurrTime.Day(), hour, minute, second, 0, timezone)
		endingHour := time.Date(maxCurrTime.Year(), maxCurrTime.Month(), maxCurrTime.Day(), hour1, minute1, second1, 0, timezone)
		//fmt.Println(openingHour, endingHour)
		totalTime += endingHour.Sub(openingHour).Hours()
		prevTimestamp := openingHour
		lastStatus := ""
		for _, status := range storeStatus {
			newTimezone := status.TimestampUTC.In(timezone)
			diff := newTimezone.Sub(prevTimestamp).Hours()

			if status.Status == "active" && diff >= 0 && newTimezone.Before(endingHour) && newTimezone.After(openingHour) {
				// fmt.Println(newTimezone, diff)
				uptime += diff
			}
			prevTimestamp = newTimezone
			lastStatus = status.Status
		}
		diff := endingHour.Sub(prevTimestamp).Hours()
		if lastStatus == "active" || lastStatus == "" && diff >= 0 {
			uptime += diff
		}

	}
	fmt.Println(totalTime, uptime)
	downtime = totalTime - uptime

	return uptime, downtime
}

func getupdateoflastday(storeID string, timezone *time.Location, currtime time.Time) (float64, float64) {

	today := time.Date(currtime.Year(), currtime.Month(), currtime.Day(), 0, 0, 0, 0, currtime.Location())
	lastDay := today.AddDate(0, 0, -1)
	newTimezone := lastDay.In(timezone)

	weekday := int(newTimezone.Weekday()) - 1
	if weekday < 0 {
		weekday = 6 // Sunday is converted to 6
	}
	storeStatus := []db.StoreData{}
	helper.DB.Where("store_id = ? AND timestamp_utc >= ? AND timestamp_utc <= ?", storeID, lastDay, today).Find(&storeStatus)
	businessHours := []db.StoreBusinessHours{}
	helper.DB.Where("store_id = ? AND day = ?", storeID, weekday).Find(&businessHours)
	uptime, downtime := getupdateoftoday(storeID, weekday, businessHours, storeStatus, timezone, currtime)

	return uptime, downtime
}

func getupdateoflastweek(storeID string, timezone *time.Location, currtime time.Time) (float64, float64) {

	today := time.Date(currtime.Year(), currtime.Month(), currtime.Day(), 0, 0, 0, 0, currtime.Location())
	total_uptime := 0.0
	total_downtime := 0.0
	for i := 1; i < 7; i++ {
		lastDay := today.AddDate(0, 0, -i)
		newTimezone := lastDay.In(timezone)

		weekday := int(newTimezone.Weekday()) - 1
		if weekday < 0 {
			weekday = 6 // Sunday is converted to 6
		}
		storeStatus := []db.StoreData{}
		helper.DB.Where("store_id = ? AND timestamp_utc >= ? AND timestamp_utc <= ?", storeID, lastDay, today).Find(&storeStatus)
		businessHours := []db.StoreBusinessHours{}
		helper.DB.Where("store_id = ? AND day = ?", storeID, weekday).Find(&businessHours)
		uptime, downtime := getupdateoftoday(storeID, weekday, businessHours, storeStatus, timezone, currtime)
		total_uptime += uptime
		total_downtime += downtime
	}

	return total_uptime, total_downtime
}
