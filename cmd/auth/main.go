package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"

	"github.com/fudio101/fube/internal/db"
	"github.com/fudio101/fube/internal/forms"
	"github.com/fudio101/fube/internal/router"
)

// @title           Golang Gin Boilerplate
// @version         1.0
// @description     A RESTful API boilerplate with Gin Framework, PostgreSQL, Redis and JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT License
// @license.url   https://github.com/fudio101/fube/blob/master/LICENSE

// @host      localhost:9000
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	//Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	if os.Getenv("ENV") == "PRODUCTION" {
		gin.SetMode(gin.ReleaseMode)
	}

	//Custom form validator
	binding.Validator = new(forms.DefaultValidator)

	//Start PostgreSQL database
	//Example: db.GetDB() - More info in the models folder
	db.Init()

	//Start Redis on database 1 - it's used to store the JWT but you can use it for anythig else
	//Example: db.GetRedis().Set(KEY, VALUE, at.Sub(now)).Err()
	db.InitRedis(1)

	// Setup router
	r := router.SetupRouter()

	port := os.Getenv("PORT")

	log.Printf("\n\n PORT: %s \n ENV: %s \n SSL: %s \n Version: %s \n\n", port, os.Getenv("ENV"), os.Getenv("SSL"), os.Getenv("API_VERSION"))

	if os.Getenv("SSL") == "TRUE" {

		//Generated using sh generate-certificate.sh
		SSLKeys := &struct {
			CERT string
			KEY  string
		}{
			CERT: "./cert/myCA.cer",
			KEY:  "./cert/myCA.key",
		}

		err := r.RunTLS(":"+port, SSLKeys.CERT, SSLKeys.KEY)
		if err != nil {
			return
		}
	} else {
		err := r.Run(":" + port)
		if err != nil {
			return
		}
	}
}
