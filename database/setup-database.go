package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/oOWinOo/justGame/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

var Db *gorm.DB

func ConnectDatabase() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s", host, port, user, password, dbname)
	newLogger := logger.New(
		log.New(os.Stdout,"\r\n",log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,			//Time manage
			LogLevel: logger.Info,				//4 types Silent Error Warn Info  [notification]
			Colorful: true,						//color
		},
	)
	var err error
	Db,err = gorm.Open(postgres.Open(dsn),&gorm.Config{
		Logger: newLogger,
	})
	if err != nil{
		panic("Failed to connect to database")
	}
	fmt.Println("Connect to database Complete")
}

func Initialize(){
	Db.AutoMigrate(models.Market{},models.User{},models.Product{},models.ProductList{})
}

func InitializeData(){
	unSellProduct := new(models.Market)
	firstMarket := new(models.Market)
	Db.Create(unSellProduct)
	Db.Create(firstMarket)

	lightBulb := new(models.ProductList)
	lightBulb.Name = "lightBulb"
	lightBulb.DefaultPrice = 20
	lightBulb.Rarity = 10000
	Db.Create(lightBulb)

	pan := new(models.ProductList)
	pan.Name = "pan"
	pan.DefaultPrice = 15
	pan.Rarity = 10000
	Db.Create(pan)

	guitar := new(models.ProductList)
	guitar.Name = "guitar"
	guitar.DefaultPrice = 20
	guitar.Rarity = 10000
	Db.Create(guitar)

	macBook := new(models.ProductList)
	macBook.Name = "macBook"
	macBook.DefaultPrice = 2000
	macBook.Rarity = 100
	Db.Create(macBook)

	mobilePhone := new(models.ProductList)
	mobilePhone.Name = "mobilePhone"
	mobilePhone.DefaultPrice = 150
	mobilePhone.Rarity = 1000
	Db.Create(mobilePhone)
}