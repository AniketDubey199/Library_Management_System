package authmiddleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/AniketDubey199/JWT_auth/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthMiddleware(database *db.MongoDB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var tokenstring string

        // 1. Prioritize getting the token from a cookie
        if cookietoken := c.Cookies("jwt"); cookietoken != "" {
            log.Warn("Token from cookie is in use")
            tokenstring = cookietoken
        } else {
            // 2. If no cookie, try to get it from the Authorization header
            authheader := c.Get("Authorization")
            if authheader == "" {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header not found"})
            }
            tokenParts := strings.Split(authheader, " ")
            if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
                return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
            }
            tokenstring = tokenParts[1]
        }
        
        // 3. Check if we have a token string to validate. If not, it's unauthorized.
        if tokenstring == "" {
             return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "No token provided"})
        }

        // 4. Perform a unified validation on the extracted token string
        secret := []byte(os.Getenv("JWT_SECRET"))
        token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
            if t.Method.Alg() != jwt.GetSigningMethod("HS256").Alg() {
                return nil, fmt.Errorf("unexpected signed method: %v", t.Header["alg"])
            }
            return secret, nil
        })

        if err != nil || !token.Valid {
            c.ClearCookie("jwt")
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
        }

        // 5. Safely extract claims using the comma-ok idiom
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
        }

        userID, ok := claims["userID"].(string)
        if !ok {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found or invalid format in token"})
        }

        if _, err := primitive.ObjectIDFromHex(userID); err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID format"})
        }

        c.Locals("userID", userID)

        return c.Next()
    }
}