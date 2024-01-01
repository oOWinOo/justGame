package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/oOWinOo/justGame/database"
	"github.com/oOWinOo/justGame/models"
)

func main() {
	database.ConnectDatabase()
	database.Initialize()
	if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }
	app := fiber.New()
	app.Static("/", "./static")
	app.Get("/login", func(c *fiber.Ctx) error {
        return c.Render("static/login.html",nil)
    })
	app.Get("/register", func(c *fiber.Ctx) error {
        return c.Render("static/register.html",nil)
    })

    app.Get("/landing", func(c *fiber.Ctx) error {
        return c.Render("static/landingPage.html", nil)
    })


	// app.Get("/",getTest)

	app.Post("/register",register)
	app.Post("/login",login)

	app.Use("/user",UserAuthRequired)
	app.Get("/user",getUserById)
	app.Get("/user/product",findProductByUserId)


	app.Get("/product",getProducts)

	app.Post("/admin/product/add",addNewProduct)
	

	app.Listen(":8080")

}

func UserAuthRequired(c *fiber.Ctx)error{
	cookie := c.Cookies("jwt")
	jwtSecretKey := os.Getenv("jwtSecretKey")
	if cookie == "" {
        return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
    }
	token,err := jwt.ParseWithClaims(cookie,jwt.MapClaims{},func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtSecretKey),nil
	})
	if err != nil || !token.Valid{
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims,ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	userID := int(claims["user_id"].(float64))
	c.Locals("user_id",userID)
	
	return c.Next()
}

func getTest(c *fiber.Ctx)error{
	fmt.Println("Get")
	return c.JSON("nothing")
}

func register(c *fiber.Ctx)error  {
	newUser := new(models.User)
	if err := c.BodyParser(newUser); err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CreateUser(database.Db,newUser);err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.SendStatus(fiber.StatusCreated)
}

func login(c *fiber.Ctx)error{
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

func getUserById(c *fiber.Ctx)error{
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

func findProductByUserId(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(user.Products)
}

func addNewProduct(c *fiber.Ctx)error{
	product := new(models.ProductList)
	if err := c.BodyParser(product);err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CreateNewProduct(database.Db,product); err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString("New Product Created")
}

func getProducts(c *fiber.Ctx)error{
	products,err := models.GetAllProducts(database.Db)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(products)
}