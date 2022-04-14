package main

import (
	"fmt"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func connect() {
	dsn := "sqlserver://KORIGOVI-PC/KorigovI:SQLEXPRESS@10.89.0.104:1433?database=Gisp"

	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("filed connection")
	}
	db.Migrator().CreateTable(&Org{})
}
