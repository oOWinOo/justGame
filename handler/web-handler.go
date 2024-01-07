package handler

import "github.com/gofiber/fiber/v2"

func LoginPage(c *fiber.Ctx) error {
	return c.Render("static/login.html", nil)
}

func RegisterPage(c *fiber.Ctx) error {
	return c.Render("static/register.html",nil)
}

func LandingPage(c *fiber.Ctx) error {
	return c.Render("static/landingPage.html", nil)
}