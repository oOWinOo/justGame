package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/oOWinOo/justGame/database"
	"github.com/oOWinOo/justGame/models"
)



func AddNewProduct(c *fiber.Ctx)error{
	product := new(models.ProductList)
	if err := c.BodyParser(product);err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CreateNewProduct(database.Db,product); err != nil{
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString("New Product Created")
}

func RandomReceiveProduct(c *fiber.Ctx)error{
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

func GetProducts(c *fiber.Ctx)error{
	products,err := models.GetAllProducts(database.Db)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	return c.JSON(products)
}

func SellProductToSystem(c *fiber.Ctx)error{
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
	if err := models.CreateHistory(database.Db,uint(userId),1,cost,productFull.Name);err!= nil{
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already Sell Product ID : %d",product.ID))

}