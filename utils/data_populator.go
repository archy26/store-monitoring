package utils

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/store-monitoring/db"
	"github.com/store-monitoring/helper"
)

func ReadCSV(fileName string) ([][]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, err
	}

	return records, nil
}

func PopulateData() error {
	//deleting previous stored data
	if helper.DB.Migrator().HasTable(&db.StoreBusinessHours{}) {
		helper.DB.Migrator().DropTable(&db.StoreBusinessHours{})
	}
	if helper.DB.Migrator().HasTable(&db.StoreData{}) {
		helper.DB.Migrator().DropTable(&db.StoreData{})
	}
	if helper.DB.Migrator().HasTable(&db.Timezone{}) {
		helper.DB.Migrator().DropTable(&db.Timezone{})

	}
	if helper.DB.Migrator().HasTable(&db.ReportStatus{}) {
		helper.DB.Migrator().DropTable(&db.ReportStatus{})

	}
	//initializing again
	err := helper.DB.AutoMigrate(&db.StoreBusinessHours{})
	if err != nil {
		return err
	}
	err = helper.DB.AutoMigrate(&db.StoreData{})
	if err != nil {
		return err
	}
	err = helper.DB.AutoMigrate(&db.Timezone{})
	if err != nil {
		return err
	}
	err = helper.DB.AutoMigrate(&db.ReportStatus{})
	if err != nil {
		return err
	}
	err = populateStoreStatus("data/store_status.csv")
	if err != nil {
		return err
	}
	err = populateTimezone("data/timezone.csv")
	if err != nil {
		return err
	}
	err = populateBusinesshours("data/business_hours.csv")
	if err != nil {
		return err
	}
	return nil
}

func populateStoreStatus(filename string) error {
	storeStatusRecords, storeStatusErr := ReadCSV(filename)
	if storeStatusErr != nil {
		return storeStatusErr
	}
	for ind, record := range storeStatusRecords {
		if ind == 0 {
			continue
		}
		storeStatusRecord := &db.StoreData{StoreID: record[0], Status: record[1], TimestampUTC: StringToTime(record[2])}
		//fmt.Println(record[0], record[1])
		helper.DB.Create(storeStatusRecord)
	}
	return nil
}
func populateTimezone(filename string) error {
	timezoneRecords, timezoneErr := ReadCSV(filename)
	if timezoneErr != nil {
		return timezoneErr
	}
	for ind, record := range timezoneRecords {
		if ind == 0 {
			continue
		}
		timezone := &db.Timezone{StoreID: record[0], TimezoneStr: record[1]}
		//fmt.Println(record[0], record[1])
		helper.DB.Create(timezone)
	}
	return nil
}

func populateBusinesshours(filename string) error {
	menuHoursRecords, menuHoursErr := ReadCSV(filename)
	if menuHoursErr != nil {
		return menuHoursErr
	}
	for ind, record := range menuHoursRecords {
		if ind == 0 {
			continue
		}
		storeBusinessHours := &db.StoreBusinessHours{StoreID: record[0], Day: int(StringToInteger(record[1])), StartTimeLocal: record[2], EndTimeLocal: record[3]}
		//fmt.Println(record[0], record[1])
		helper.DB.Create(storeBusinessHours)
	}
	return nil
}
