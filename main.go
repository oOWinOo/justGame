package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
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
	// database.InitializeData()
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
	app.Get("/user/product/market",getUserMarketProducts)
	app.Patch("/user/product/sellsystem",sellProductToSystem)
	app.Patch("/user/product/recieve/random",randomReceiveProduct)

	app.Get("/product/market/:id",getMarketProducts)
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

func getUserMarketProducts(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	marketProduct := new([]models.Product)
	for _,product := range user.Products{
		if product.MarketId != 1{
			*marketProduct = append(*marketProduct, product)
		}
	}
	fmt.Println(*marketProduct)
	return c.JSON(*marketProduct)
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

func randomReceiveProduct(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user.Money < 50{
		return c.Status(fiber.StatusBadRequest).SendString("Not enough money.")
	}
	err = models.ReceiveNewRandomProduct(database.Db,&user)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString("Random already")
}

func getProducts(c *fiber.Ctx)error{
	products,err := models.GetAllProducts(database.Db)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(products)
}

func getMarketProducts(c *fiber.Ctx)error{
	marketId,err := strconv.Atoi(c.Params("id"))
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}
	products,err := models.GetAllMarketProducts(database.Db,uint(marketId))
	if err != nil{
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(products)
}

func sellProductToSystem(c *fiber.Ctx)error{
	userId,ok := c.Locals("user_id").(int)

	if !ok{
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user,err := models.GetUser(database.Db,userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	
	product := new(models.Product)

	if err := c.BodyParser(product); err != nil{
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	productFull ,err := models.GetProductById(database.Db,product.ID)
	cost :=  productFull.DefaultPrice
	result := database.Db.Delete(&productFull)
	if result.Error != nil{
		return c.Status(fiber.StatusBadRequest).SendString(result.Error.Error())
	}
	user.Money += cost
	if err := database.Db.Save(user).Error; err != nil {
        return c.Status(fiber.StatusBadRequest).SendString(err.Error())
    }
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already Sell Product ID : %d",product.ID))
	
	
}