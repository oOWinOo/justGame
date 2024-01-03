package models

import "gorm.io/gorm"

type Market struct {
	gorm.Model
	Products		[]Product	`json:"market_product"`
}

func GetAllMarketProducts(Db *gorm.DB,id uint) ([]Product,error) {
	market := new(Market)
	result := Db.Preload("Products").Where("id = ?",id).First(market)
	if result.Error != nil{
		return market.Products,result.Error
	}
	return market.Products,nil
}

func GetMarketProductByPriceFilter(Db *gorm.DB,first uint,last uint) ([]Product,error) {
	market := new(Market)
	result := Db.Find(market)
	var filterProducts []Product
	if result.Error != nil{
		return filterProducts,result.Error
	}
	for _,product := range market.Products{
		if product.DefaultPrice >= first && product.DefaultPrice <= last{
			filterProducts = append(filterProducts,product)
		}
	}

	return filterProducts,nil
}