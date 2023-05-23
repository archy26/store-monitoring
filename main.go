package main

import (
	"github.com/store-monitoring/helper"
	"github.com/store-monitoring/server"
	"github.com/store-monitoring/utils"
)

func main() {
	helper.Init()
	err := utils.PopulateData()
	if err != nil {
		panic(err)
	}
	server.ServeRequest()

}
