package main

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type DbContext struct {
	db *gorm.DB
}

func connect(db *gorm.DB) {
	dsn := "sqlserver://@localhost:52876?database=Gisp"

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("filed connection")
	} else {
		fmt.Println("Success connect to MSSQL")
	}

}

func getAll(db *gorm.DB) []Org {
	var result []Org
	db.Find(&result)
	return result
}

func createOrg(db *gorm.DB, data *[]Org) {
	db.Create(data)
}

func createProds(db *gorm.DB, data *[]Prod) {
	db.Create(data)
}
