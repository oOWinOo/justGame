package models

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name		string				`json:"user_name" gorm:"unique"`
    Email    	string 				`json:"email" gorm:"unique"`
	Password	string				`json:"password"`
	Money		uint				`json:"money"`
	Products	[]Product			`json:"products"`
	StorageSize	uint				`json:"storage"`
	LastLogin	time.Time			`json:"lastlogin"`
}	

func CreateUser(Db *gorm.DB,user *User) error{

	if !isValidEmail(user.Email) {
        return errors.New("Invalid email format")
    }

	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	user.Password = string(hashedPassword)
	user.Money = 1000
	user.StorageSize = 10
	result := Db.Create(user)
	if result.Error != nil{
		return result.Error
	}
	err = ReceiveNewRandomProduct(Db,user)
	if err != nil{
		return err
	}
	err = ReceiveNewRandomProduct(Db,user)
	if err != nil{
		return err
	}
	return nil
}

func LoginUser(Db *gorm.DB,user *User) (string,error){
	searchUser := new(User)
	result := Db.Where("email = ?",user.Email).First(searchUser)
	if result.Error != nil{
		return "",result.Error
	}
	if !isGetDaily(searchUser){
		if err := getDailyMoney(Db,searchUser); err!= nil{
			return "",err
		}
	}
	if err := UpdateLastLogin(Db,searchUser); err!= nil{
		return "",err
	}
	err := bcrypt.CompareHashAndPassword([]byte(searchUser.Password),[]byte(user.Password))
	if err != nil{
		return "",err
	}


    if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file")
        return "",err
    }

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)// ใส่ข้อมูลอื่นๆ เข้าไป  claims
	claims["user_id"] = searchUser.ID
	claims["exp"] = time.Now().Add(time.Hour*72).Unix()

	jwtSecretKey :=  os.Getenv("jwtSecretKey")
	tokenString,err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "",err
	}
	
	return tokenString,nil
}

func GetUser(Db *gorm.DB,id int)(User,error){
	user := new(User)
	result := Db.Preload("Products").Where("id = ?",id).First(user)
	if result.Error != nil{
		return *user,result.Error
	}
	return *user,nil
}

func GetUserProducts(Db *gorm.DB,userId int)([]Product,error){
	user,err := GetUser(Db,userId)
	if err != nil {
		return nil,err
	}
	return user.Products,nil
}


func ReceiveNewRandomProduct(Db *gorm.DB,user *User)error{
	if !CheckStorage(user.StorageSize,user.Products){
		return errors.New("Not enough storage.")
	}
	product,err := RandomProduct(Db)
	if err != nil{
		return err
	}
	product.UserId = user.ID
	result := Db.Create(product)
	if result.Error != nil{
		return result.Error
	}

	user.Products = append(user.Products, *product)
	user.Money = user.Money - 50
	if err := Db.Save(user).Error; err != nil {
        return err
    }
	if err := CreateHistory(Db,user.ID,4,50,product.Name);err!= nil{
		return err
	}
	return nil
}

func isValidEmail(email string) bool {
    pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    regex := regexp.MustCompile(pattern)
	
    return regex.MatchString(email)
}

func CheckStorage(cap uint,products []Product) bool{
	for _,product := range products{
		if !product.Sold {
			cap -= 1
		}
	}
	if cap <= 0 {
		return false
	}
	return true
}

func UpdateLastLogin(Db *gorm.DB,user *User)error{
	user.LastLogin = time.Now()
	if err := Db.Save(user).Error; err != nil{
		return err
	}
	return nil

}

func isGetDaily(user *User)bool{
	lastLogin := user.LastLogin
	now := time.Now()
	return now.Truncate(12*time.Hour).Equal(lastLogin.Truncate(12*time.Hour))
}

func getDailyMoney(Db *gorm.DB,user *User)error{
	dailyMoney,err := CalDailyMoney(Db,user)
	if err != nil{
		return err
	}
	user.Money += dailyMoney
	if err := Db.Save(user).Error; err != nil{
		return err
	}
	if err := CreateHistory(Db,user.ID,5,dailyMoney,"");err!= nil{
		return err
	}
	return nil
}

func CalDailyMoney(Db *gorm.DB,user *User)(uint,error){
	initialMoney := 100
	products,err := GetUserProducts(Db,int(user.ID))
	if err != nil {
		return 0,err
	}
	for _,product := range products{
		if product.MarketId != 1{
			continue
		}
		rarity ,err := findRarity(Db,product.Name)
		if err != nil{
			return 0,err
		}
		fmt.Println(rarity)
		fmt.Println(initialMoney)
		initialMoney += 10000/int(rarity)
	}
	return uint(initialMoney),nil
}

func findRarity(Db *gorm.DB,name string)(uint ,error){
	productList := new(ProductList)
	result := Db.Where("name = ?",name).First(productList)
	if result.Error != nil {
		return 0,result.Error
	}
	return productList.Rarity,nil

}

func Upgrade(Db *gorm.DB,user *User)error{
	cost := user.StorageSize*10
	if user.Money < cost{
		return errors.New("Not enough money.")
	}
	user.Money -= cost
	user.StorageSize += 10;
	if err := Db.Save(user).Error; err != nil{
		return err
	}
	if err := CreateHistory(Db,user.ID,6,cost,"");err!= nil{
		return err
	}
	return nil
}