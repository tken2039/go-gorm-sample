package main

import (
	"encoding/json"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Fruit struct {
	FruitID    uint     `gorm:"column:fruit_id"`
	MarketID   uint     `gorm:"column:market_id"`
	Fruit      string   `gorm:"column:fruit_name"`
	CustomerID uint     `gorm:"column:customer_id"`
	Customer   Customer `gorm:"foreignKey:CustomerID;references:CustomerID;"` // Define external references on gorm.
}

// Table names referenced by gorm are by default of the form "StructName" + "s".
// If the table name is not in the above format, it is necessary to fix the table name referred
// to by gorm by implementing the TableName method as follows.
func (f *Fruit) TableName() string {
	return "fruit"
}

type Customer struct {
	CustomerID   uint   `gorm:"column:customer_id"`
	CustomerName string `gorm:"column:customer_name"`
}

func (c *Customer) TableName() string {
	return "customer"
}

type Market struct {
	MarketID   uint    `gorm:"column:market_id"`
	MarketName string  `gorm:"column:market_name"`
	Fruits     []Fruit `gorm:"foreignKey:MarketID;references:MarketID;"` // Define external references on gorm.
}

func main() {
	// connect to DB
	dsn := "root:password@tcp(127.0.0.1:3306)/go_gorm_sample?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// find data
	var data []Market
	db.Table("market"). // Synonymous with Model(&Market{}). But in that case, gorm tries to refer to the markets table, so the TableName method must be defined.
				Preload("Fruits.Customer"). // fruit has single customer (`has one` pattern). In the case of `has one`, do not put `s` after the name of the structure. If the structure is nested, then `. ` to represent nested structures.
				Preload("Fruits").          // market has multiple fruits (`has many` pattern). In the case of `has many`, add `s` after the name of the structure.
				Find(&data)

	// Convert the output to json for easier viewing.
	b, _ := json.Marshal(data)

	fmt.Println(string(b))
}
