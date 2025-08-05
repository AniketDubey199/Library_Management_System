package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/AniketDubey199/JWT_auth/db"
)

func Authentication(router fiber.Router , database *MongoDB){
	router.Post("/register" , func (c *fiber.Ctx) error {
		user := &User{
			Username : c.FormValue("username"),
    		Password : c.FormValue("password"),
		}

		if ( user.Username == "" || user.Password== ""){
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Username and Password cannot be empty",
			})
		}

		hashed , err := bcrypt.GenerateFromPassword([]byte(user.Password) , bcrypt.DefaultCost)
		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":err.Error(),
			})
		}

		user.Password = string(hashed);
		db.CreateUser(user)

		
	})
	router.Post("/login" , func (c *fiber.Ctx) error {
		return c.SendString("login")
	})
}