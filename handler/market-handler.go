package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/oOWinOo/justGame/database"
	"github.com/oOWinOo/justGame/models"
)

func GetUserMarketProducts(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(int)

	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	user, err := models.GetUser(database.Db, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	marketProduct := new([]models.Product)
	for _, product := range user.Products {
		if product.MarketId != 1 {
			*marketProduct = append(*marketProduct, product)
		}
	}
	return c.JSON(*marketProduct)
}

func GetMarketProducts(c *fiber.Ctx) error {
	marketId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	products, err := models.GetAllMarketProducts(database.Db, uint(marketId))
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(products)
}

func SellProductToMarket(c *fiber.Ctx) error {
	type ProductRequest struct {
		ID       uint   `json:"ID"`
		PriceSet string `json:"price_set"`
	}
	product := new(ProductRequest)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	productFull, err := models.GetProductById(database.Db, product.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	cost, err := strconv.ParseUint(product.PriceSet, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.PostProductOnMarket(database.Db, &productFull, uint(cost)); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already post Product ID : %d from Market", product.ID))
}

func BuyProductFromMarket(c *fiber.Ctx) error {
	type ProductRequest struct {
		ID       uint   `json:"ID"`
		PriceSet string `json:"price_set"`
	}
	product := new(ProductRequest)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	beforeProduct, err := models.GetProductById(database.Db, product.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	userId, ok := c.Locals("user_id").(int)

	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user, err := models.GetUser(database.Db, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if user.Money < beforeProduct.PriceSet {
		return c.Status(fiber.StatusBadRequest).SendString("Not enough money.")
	}
	if !models.CheckStorage(user.StorageSize, user.Products) {
		return c.Status(fiber.StatusBadRequest).SendString("Not enough storage.")
	}
	sellerId := beforeProduct.UserId
	buyerId := user.ID
	cost := beforeProduct.PriceSet

	if err := models.ChangeOwner(database.Db, user.ID, &beforeProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	buyUser, err := models.GetUser(database.Db, int(buyerId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	buyUser.Money -= cost
	if err := database.Db.Save(buyUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	sellUser, err := models.GetUser(database.Db, int(sellerId))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	sellUser.Money += cost
	if err := database.Db.Save(sellUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CreateHistory(database.Db, sellerId, 2, cost, beforeProduct.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	if err := models.CreateHistory(database.Db, buyerId, 3, cost, beforeProduct.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already buy Product ID : %d from Market", product.ID))
}

func CancelProductFromMarket(c *fiber.Ctx) error {
	userId, ok := c.Locals("user_id").(int)

	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	user, err := models.GetUser(database.Db, userId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if !models.CheckStorage(user.StorageSize, user.Products) {
		return c.Status(fiber.StatusBadRequest).SendString("Not enough storage.")
	}

	type ProductRequest struct {
		ID uint `json:"ID"`
	}
	product := new(ProductRequest)
	if err := c.BodyParser(product); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	productFull, err := models.GetProductById(database.Db, product.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := models.CancelProductOnMarket(database.Db, &productFull); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("Already cancel sell Product ID : %d from Market", product.ID))

}