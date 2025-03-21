package users

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	// "golang.org/x/crypto/bcrypt"

	"github.com/adarsh-jaiss/zocket/internal/middleware"
	"github.com/adarsh-jaiss/zocket/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// SignupRequest represents the request body for user signup
type SignupRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Signup handles user registration and returns a JWT token
func Signup(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req SignupRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Hash the password
		// hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		// 		"error": "Failed to process password",
		// 	})
		// }

		// Create user object
		user := types.User{
			Email:     req.Email,
			Password:  req.Password,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		// Store user in database
		userID, err := CreateUserInStore(db, user)
		if err != nil {
			if err == ErrEmailExists {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "Email already exists",
				})
			}
			fmt.Println(err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				
				"error": "Failed to create user",
			})
		}

		// Generate JWT token
		claims := jwt.MapClaims{
			"user_id": userID,
			"email":   user.Email,
			"exp":     time.Now().Add(time.Hour * time.Duration(middleware.JWTExpirationHours)).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(middleware.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"user_id": userID,
			"token":   t,
			"message": "User created successfully",
		})
	}
}

// GetUser handles the GET request for retrieving a user
func GetUser(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		// Get user claims from JWT token
		token := c.Locals("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		tokenUserID := int(claims["user_id"].(float64))

		// Only allow users to access their own data
		if userID != tokenUserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}

		user, err := GetUserFromStore(db, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve user",
			})
		}

		// Don't return password in response
		user.Password = ""
		return c.Status(fiber.StatusOK).JSON(user)
	}
}




func SignIn(db *sql.DB) fiber.Handler {
	return func (c *fiber.Ctx) error  {
		var req types.SignInRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		user, err := GetUserByEmailAndPassword(db, req.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Invalid email",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to sign in",
			})
		}

		// compare the password
		if user.Password != req.Password {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid password",
			})
		}

		// Generate JWT token
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email":   user.Email,
			"exp":     time.Now().Add(time.Hour * time.Duration(middleware.JWTExpirationHours)).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(middleware.JWTSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate token",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"user_id": user.ID,
			"token":   t,
			"message": "User signed in successfully",
		})
	}
}

func FetchAllUsers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := GetAllUsers(db)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to retrieve users",
			})
		}

		return c.Status(fiber.StatusOK).JSON(users)
	}
}