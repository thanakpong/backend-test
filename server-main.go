package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	models "github.com/thanakpong/test-api-demo/data-models"
	"github.com/thanakpong/test-api-demo/storage"

	"gorm.io/gorm"
)

type Post struct {
	Id         uint      `gorm:"primary key;autoIncrement" json:"id"`
	Title      string    `json:"Title"`
	Content    string    `json:"content"`
	Published  bool      `json:"Published"`
	View_count int       `json:"View_count"`
	Created_at time.Time `json:"Created_at"`
	Update_at  time.Time `json:"Update_at"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreatePost(context *fiber.Ctx) error {
	post := Post{}

	err := context.BodyParser(&post)

	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"massage": "request failed"})
		return err
	}

	err = r.DB.Create(&post).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"massage": "could not create post"})
		return err
	}
	context.Status(http.StatusBadRequest).JSON(&fiber.Map{
		"massage": "post has been added"})
	return nil

}

func (r *Repository) DeletePost(context *fiber.Ctx) error {
	postModel := &[]models.Post{}
	id := context.Params("id")
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"massage": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(postModel, id)

	if err.Error != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"massage": "cou;d not delete",
		})
		return err.Error
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"massage": "post delete successfully",
	})
	return nil
}

func (r *Repository) GetPost(context *fiber.Ctx) error {
	postModels := &[]models.Post{}

	err := r.DB.Find(postModels).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"massage": "could not get post"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"massage": "post fetched successfully",
		"data":    postModels})
	return nil
}

func (r *Repository) GetPostID(context *fiber.Ctx) error {
	id := context.Params("id")
	postModel := &models.Post{}
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"massage": "id cannot be empty",
		})
		return nil
	}
	fmt.Println("Id is", id)

	err := r.DB.Where("id = ?", id).First(postModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"massage": "could not get post",
		})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"massage": "post is fetched sucsessfully",
	})
	return nil

}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_post", r.CreatePost)
	api.Delete("delete_post/:id", r.DeletePost)
	api.Get("/post", r.GetPost)
	api.Get("/get_post/:id", r.GetPostID)
	// api.Put()

}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("could not load the detabase")
	}
	err = models.MigratePost(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8000")

}
