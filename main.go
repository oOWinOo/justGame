package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/oOWinOo/justGame/database"
	"github.com/oOWinOo/justGame/handler"
)

func main() {
	database.ConnectDatabase()
	database.Initialize()
	database.InitializeData()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := fiber.New()
	app.Static("/", "./static")
	app.Get("/", handler.LoginPage)
	app.Get("/register", handler.RegisterPage)
	app.Get("/landing", handler.LandingPage)

	// app.Post("/",postTest)
	app.Post("/register", handler.Register)
	app.Post("/login", handler.Login)

	app.Use("/user", handler.UserAuthRequired)
	app.Get("/user", handler.GetUserById)
	app.Get("/user/loginReward", handler.GetUserLoginReward)
	app.Get("/user/product", handler.GetUserProducts)
	app.Get("/user/product/inventory", handler.GetUserInventoryProducts)
	app.Get("/user/product/market", handler.GetUserMarketProducts)
	app.Get("/user/product/history", handler.GetHistory)
	app.Patch("/user/product/sellsystem", handler.SellProductToSystem)
	app.Patch("/user/product/sellmarket", handler.SellProductToMarket)
	app.Patch("/user/product/buymarket", handler.BuyProductFromMarket)
	app.Patch("/user/product/cancelmarket", handler.CancelProductFromMarket)
	// app.Patch("/user/product/claim",claimMoneyFromSoldProduct)
	app.Patch("/user/product/recieve/random", handler.RandomReceiveProduct)
	app.Patch("/user/product/upgradeinventory", handler.UpgradeInventory)

	app.Get("/product/market/:id", handler.GetMarketProducts)
	app.Get("/product", handler.GetProducts)

	app.Post("/admin/product/add", handler.AddNewProduct)

	app.Listen(":8080")

}

// func postTest(c *fiber.Ctx)error{
// 	type ProductRequest struct {
// 		ID        uint   `json:"ID"`
// 	}
// 	product := new(ProductRequest)
// 	if err := c.BodyParser(product); err != nil{
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}
// 	beforeProduct ,err := models.GetProductById(database.Db,product.ID)
// 	if err != nil{
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}
// 	var p []uint
// 	p = append(p, beforeProduct.UserId)
// 	models.ChangeOwner(database.Db,5,&beforeProduct)

// 	afterProduct ,err := models.GetProductById(database.Db,product.ID)
// 	if err != nil{
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}

// 	p = append(p, afterProduct.UserId)
// 	return c.JSON(p)
// }

// func claimMoneyFromSoldProduct(c *fiber.Ctx)error{
// 	userId,ok := c.Locals("user_id").(int)
// 	if !ok{
// 		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
// 	}
// 	user,err := models.GetUser(database.Db,userId)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}
// 	product := new(models.Product)

// 	if err := c.BodyParser(product); err != nil{
// 		fmt.Println(err.Error())
// 		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
// 	}
// 	productFull ,err := models.GetProductById(database.Db,product.ID)
// 	cost :=  productFull.PriceSet
// 	result := database.Db.Delete(&productFull)
// 	if result.Error != nil{
// 		return c.Status(fiber.StatusBadRequest).SendString(result.Error.Error())
// 	}
// 	user.Money += cost
// 	if err := database.Db.Save(user).Error; err != nil {
//         return c.Status(fiber.StatusBadRequest).SendString(err.Error())
//     }
// 	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Receive : %d from sell product : %s in the market",productFull.PriceSet,productFull.Name))
// }
