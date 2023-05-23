package main

import (
	"github.com/store-monitoring/helper"
	"github.com/store-monitoring/server"
)

func main() {
	helper.Init()
	// err := utils.PopulateData()
	// if err != nil {
	// 	panic(err)
	// }
	server.ServeRequest()

}
