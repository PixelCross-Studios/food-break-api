package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	adapter "github.com/gwatts/gin-adapter"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Product struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Code  string  `json:"code"`
	Price float32 `json:"price"`
}

func SetupRouter() *gin.Engine {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	issuerURL, _ := url.Parse(os.Getenv("AUTH0_ISSUER_URL"))
	audience := os.Getenv("AUTH0_AUDIENCE")

	provider := jwks.NewCachingProvider(issuerURL, time.Duration(5*time.Minute))

	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
	)

	jwtMiddleware := jwtmiddleware.New(jwtValidator.ValidateToken)
	r.Use(adapter.Wrap(jwtMiddleware.CheckJWT))

	r.GET("/products", func(c *gin.Context) {
		products := []Product{
			{ID: 1, Title: "Product 1", Code: "p1", Price: 100.0},
			{ID: 2, Title: "Product 2", Code: "p2", Price: 200.0},
			{ID: 3, Title: "Product 3", Code: "p3", Price: 300.0},
		}
		c.JSON(http.StatusOK, products)
	})

	return r
}

func main() {
	r := SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
