package models

import (
	"errors"

	"math/rand"

	"gorm.io/gorm"
)

type ProductList struct {
	Name			string	`json:"product_name" gorm:"unique"`
	DefaultPrice	uint	`json:"default_price"`
	Rarity			uint	`json:"rarity"` // Most rare = 1
}

type Product struct{
	gorm.Model
	Name			string	`json:"product_name"`
	DefaultPrice	uint	`json:"default_price"`
	UserId			uint	`json:"user_id"`
}

func CreateNewProduct(Db *gorm.DB,product *ProductList)error{
	if product.DefaultPrice == 0{
		return errors.New("Default Price can not be 0")
	}
	if product.Rarity == 0{
		return errors.New("Rarity can not be 0")
	}
	result := Db.Create(product)
	if result.Error != nil{
		return result.Error
	}
	return nil
}

func RandomProduct(Db *gorm.DB)(*Product,error){
	var productList []ProductList
	newProduct := new(Product)
	result := Db.Find(&productList)
	if result.Error != nil{
		return newProduct,result.Error
	}
	var totalRate uint
	for _,product := range productList{
		totalRate += product.Rarity
	}
	if totalRate == 0 {
		return newProduct, errors.New("No products available")
	}
	randomNumber := rand.Intn(int(totalRate))
	var findingRate uint
	for _,product := range productList{
		findingRate += product.Rarity
		if findingRate > uint(randomNumber){

			newProduct.Name = product.Name
			newProduct.DefaultPrice = product.DefaultPrice
			return newProduct,nil
		}
	}
	return newProduct,errors.New("Can't Random Product")

}

func GetAllProducts(Db *gorm.DB)([]Product,error){
	products := new([]Product)
	result := Db.Find(&products)
	if result.Error != nil{
		return *products,result.Error
	}
	return *products,nil
}

