package book

import (
	// "github.com/AniketDubey199/JWT_auth/book"
	// "strconv"

	"context"
	"strconv"

	// "github.com/AniketDubey199/JWT_auth/book"
	"github.com/AniketDubey199/JWT_auth/db"
	"github.com/AniketDubey199/JWT_auth/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

func Bookhandler(route fiber.Router, database *db.MongoDB) {
	route.Get("/", func(c *fiber.Ctx) error {
		userIDInterface := c.Locals("userID")

		userID := userIDInterface.(string)

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert it into primitive object id",
			})
		}

		filter := bson.M{"userid": objID}

		if title := c.Query("title"); title != "" {
			filter["title"] = bson.M{"$regex": primitive.Regex{Pattern: title, Options: "i"}}
		}
		if author := c.Query("author"); author != "" {
			filter["author"] = bson.M{"$regex": primitive.Regex{Pattern: author, Options: "i"}}
		}
		if status := c.Query("status"); status != "" {
			filter["status"] = status
		}
		if yearStr := c.Query("year"); yearStr != "" {
			if year, err := strconv.Atoi(yearStr); err == nil {
				filter["year"] = year
			}
		}

		var book []model.Book
		collection := database.Client.Database("librarydb").Collection("book")
		cursor, err := collection.Find(context.Background(), filter)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		defer cursor.Close(context.Background())

		if err := cursor.All(context.Background(), &book); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)

	})
	route.Get("/:id", func(c *fiber.Ctx) error {
		bookIDHex := c.Params("id")
		bookID, err := primitive.ObjectIDFromHex(bookIDHex)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID format."})
		}
		userIDInterface := c.Locals("userID")

		userID := userIDInterface.(string)

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert it into primitive object id",
			})
		}

		var book model.Book
		filter := bson.M{"_id": bookID, "userid": objID}

		booksCollection := database.Client.Database("librarydb").Collection("book")
		none := booksCollection.FindOne(context.Background(), filter).Decode(&book)
		if none != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cannot found the book in database",
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)

	})
	route.Post("/", func(c *fiber.Ctx) error {
		book := new(model.Book)
		userIDInterface := c.Locals("userID")

		userID := userIDInterface.(string)

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert it into primitive object id",
			})
		}

		book.UserID = objID

		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		// dont forget to add the validator

		if err := db.AddBook(database.Client, book); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(book)
	})
	route.Put("/:id", func(c *fiber.Ctx) error {
		bookIDHex := c.Params("id")
		bookID, err := primitive.ObjectIDFromHex(bookIDHex)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID format."})
		}
		userIDInterface := c.Locals("userID")

		userID := userIDInterface.(string)

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert it into primitive object id",
			})
		}

		var updateData struct {
			Title  string `bson:"title,omitempty"`
			Status string `bson:"status,omitempty"`
			Author string `bson:"author,omitempty"`
			Year   int    `bson:"year,omitempty"`
		}

		if err := c.BodyParser(&updateData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		filter := bson.M{"_id": bookID, "userid": objID}

		update := bson.M{"$set": updateData}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var book model.Book
		booksCollection := database.Client.Database("librarydb").Collection("book")
		none := booksCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&book)
		if none != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Cannot found the book in database",
			})
		}
		return c.Status(fiber.StatusOK).JSON(book)
	})
	route.Delete("/:id", func(c *fiber.Ctx) error {
		bookIDHex := c.Params("id")
		bookID, err := primitive.ObjectIDFromHex(bookIDHex)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid book ID format."})
		}
		userIDInterface := c.Locals("userID")

		userID := userIDInterface.(string)

		objID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot convert it into primitive object id",
			})
		}

		// var book model.Book
		filter := bson.M{"_id": bookID, "userid": objID}
		booksCollection:= database.Client.Database("librarydb").Collection("book")
		none := booksCollection.FindOneAndDelete(context.Background() , filter).Err()
		if none != nil {
			if none == mongo.ErrNoDocuments {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "Book not found or does not belong to the user.",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": none.Error(),
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":"Book deleted succesfully",
		})

	})

}
