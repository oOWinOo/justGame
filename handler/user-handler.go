package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/oOWinOo/justGame/database"
	"github.com/oOWinOo/justGame/models"
)

func UserAuthRequired(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	jwtSecretKey := os.Getenv("jwtSecretKey")
	if cookie == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}
	token, err := jwt.ParseWithClaims(cookie, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID := int(claims["user_id"].(float64))
	c.Locals("user_id", userID)

	return c.Next()
}

func Register(c *fiber.Ctx)error  {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CreateUser(database.Db,newUser);err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

func Login(c *fiber.Ctx)error{
	user := new(models.User)
	if err := c.BodyParser(user); err !=nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	token,err := models.LoginUser(database.Db,user)
	if err != nil{
		return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:"jwt",
		Value:token,
		Expires: time.Now().Add(time.Hour*72),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).SendString("Login Success")
}

func GetUserById(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)
	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user)
}

func GetUserLoginReward(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)
	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	reward,err := models.CalDailyMoney(database.Db,&user)
	return c.JSON(reward)
}

func GetUserProducts(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	products,err := models.GetUserProducts(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(products)
}

func GetUserInventoryProducts(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)
	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	var notInMarket []models.Product
	for _,product := range user.Products{
		if product.MarketId == 1{
			notInMarket = append(notInMarket, product)
		}
	}
	return c.JSON(notInMarket)
}

func GetHistory(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	historys,err := models.GetHistoryByUserId(database.Db,uint(userId))
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).SendString("Can't get history")
	}
	return c.JSON(historys)
}

func UpgradeInventory(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.Upgrade(database.Db,&user) ; err!=nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already upgrade inventory"))
}