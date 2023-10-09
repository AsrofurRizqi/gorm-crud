package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AsrofurRizqi/gorm-crud/models"
	"github.com/AsrofurRizqi/gorm-crud/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	City  string `json:"city"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetUsers(c *fiber.Ctx) error {
	userModels := &[]models.User{}

	find := r.DB.Find(&userModels).Error

	if find != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "could not get users",
				"error":   find,
			},
		)
		return find
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "success",
			"users":   userModels,
		},
	)
	return nil
}

func (r *Repository) GetUser(c *fiber.Ctx) error {
	userModels := &models.User{}

	id := c.Params("id")

	//check id empty or not integer
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is empty",
			},
		)
		return nil
	} else if _, err := strconv.Atoi(id); err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is not valid",
			},
		)
		return nil
	}

	find := r.DB.First(&userModels, id).Error

	if find == gorm.ErrRecordNotFound {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "could not get user, user not found",
				"status":  404,
			},
		)
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "success",
			"user":    userModels,
		},
	)
	return nil
}

func (r *Repository) NewUser(c *fiber.Ctx) error {
	user := User{}

	err := c.BodyParser(&user)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "body parsing failed",
				"error":   err,
			},
		)
		return err
	}

	add := r.DB.Create(&user).Error
	if add != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "insert database failed",
				"error":   add,
			},
		)
		return add
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "User created successfully",
			"user":    user,
		},
	)
	return nil
}

func (r *Repository) DeleteUser(c *fiber.Ctx) error {
	user := User{}

	id := c.Params("id")

	//check id empty or not integer
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is empty",
			},
		)
		return nil
	} else if _, err := strconv.Atoi(id); err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is not valid",
			},
		)
		return nil
	}

	find := r.DB.Where("id = ?", id).First(&user).Error
	if find == gorm.ErrRecordNotFound {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "user not found",
				"status":  404,
			},
		)
	}

	delete := r.DB.Where("id = ?", id).Delete(&user).Error

	if delete != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "could not delete user",
				"error":   delete,
			},
		)
		return delete
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "success delete user",
			"user":    user,
		},
	)
	return nil
}

func (r *Repository) UpdateUser(c *fiber.Ctx) error {
	user := User{}

	id := c.Params("id")

	//check id empty or not integer
	if id == "" {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is empty",
			},
		)
		return nil
	} else if _, err := strconv.Atoi(id); err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"message": "id is not valid",
			},
		)
		return nil
	}

	err := r.DB.Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "user not found",
				"status":  404,
			},
		)
	}

	parse := c.BodyParser(&user)

	if parse != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{
				"message": "body parsing failed",
				"error":   err,
			},
		)
		return err
	}

	err = r.DB.Where("id = ?", id).Updates(&user).Error
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{
				"message": "update database failed",
				"error":   err,
			},
		)
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "User updated successfully",
			"user":    user,
		},
	)
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/users", r.GetUsers)
	api.Get("/user/:id", r.GetUser)
	api.Post("/user", r.NewUser)
	api.Delete("/user/:id", r.DeleteUser)
	api.Put("/user/:id", r.UpdateUser)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateUsers(db)
	if err != nil {
		log.Fatal("Error migrating users")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":3000")
}
