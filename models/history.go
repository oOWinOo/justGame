package models

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	UserId		uint	`json:"user_id"`
	TypeHistory	uint	`json:"is_receive"`
	Value		uint	`json:"value"`
	ProductName	string	`json:"product_name"`
}

// 1.systemsell			product  + value
// 2.marketsell			product	 + value
// 3.buymarket			product  - value
// 4.randomproduct		product	 - value
// 5.loginReward

func GetHistoryByUserId(Db *gorm.DB,id uint)([]History,error){
	historys := new([]History)
	result := Db.Where("user_id = ?",id).Order("created_at desc").Limit(20).Find(historys)
	if result.Error != nil{
		return *new([]History),result.Error
	}
	return *historys,nil
}

func CreateHistory(Db *gorm.DB,id uint,historyType uint,value uint,productName string)(error){
	history := new(History)
	history.UserId = id;
	history.TypeHistory = historyType
	history.Value = value
	history.ProductName = productName
	result := Db.Create(history)
	if result.Error != nil{
		return result.Error
	}
	return nil
}



