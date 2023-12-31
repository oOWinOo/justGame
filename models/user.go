package models
import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"


	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name		string				`json:"user_name"`
    Email    	string 				`json:"email" gorm:"unique"`
	Password	string				`json:"password"`
	Money		uint				`json:"money"`
	Products	[]Product			`json:"products"`
}	

func CreateUser(Db *gorm.DB,user *User) error{
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(user.Password),bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	user.Password = string(hashedPassword)
	user.Money = 1000
	
	result := Db.Create(user)
	if result.Error != nil{
		return result.Error
	}
	err = ReceiveNewProduct(Db,user)
	if err != nil{
		return err
	}
	err = ReceiveNewProduct(Db,user)
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


func ReceiveNewProduct(Db *gorm.DB,user *User)error{
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
	if err := Db.Save(user).Error; err != nil {
        return err
    }
	return nil
}