package auth

import (
	"context"
	// "fmt"
	"log"
	"github.com/AniketDubey199/JWT_auth/db"
	"github.com/AniketDubey199/JWT_auth/model"
	"github.com/AniketDubey199/JWT_auth/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Authentication(router fiber.Router , database *db.MongoDB){
	router.Post("/register" , func (c *fiber.Ctx) error {
		user := &model.User{
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

		user.Password = string(hashed)
		db.CreateUser(database.Client , user)

		token , err := utils.GenerateToken(user)

		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name: "jwt",
			Value: token,
			HTTPOnly: !c.IsFromLocal(),
			Secure: !c.IsFromLocal(),
			MaxAge: 3600 * 24 * 7,
		})
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token":token,
		})
		
	})
	router.Post("/login" , func (c *fiber.Ctx) error {
		dbuser := new(model.User)
		authuser := &model.User{
			Username : c.FormValue("username"),
    		Password : c.FormValue("password"),
		}

		if ( authuser.Username == "" || authuser.Password== ""){
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":"Username and Password cannot be empty",
			})
		}

		filter := bson.M{"username" : authuser.Username}

		err := database.Collection.FindOne(context.Background() , filter).Decode(&dbuser)

		if err != nil{
			if err == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":"Username not found",
				})
			}else{
				log.Fatal(err)
			}
		}

		if err := bcrypt.CompareHashAndPassword([]byte(dbuser.Password) , []byte(authuser.Password)) ; err != nil{
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":"Invalid Credentials",
			})
		}

		token , err := utils.GenerateToken(authuser)
		if err != nil{
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":err.Error(),
			})
		}
		c.Cookie(&fiber.Cookie{
			Name: "jwt",
			Value: token,
			HTTPOnly: !c.IsFromLocal(),
			Secure: !c.IsFromLocal(),
			MaxAge: 3600 * 24 * 7,
		})
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token":token,
		})
	})
}