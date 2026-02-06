package main

import (
	"RESTAPI/internal/api/handler"
	mid "RESTAPI/internal/api/middlerwares"
	"RESTAPI/internal/api/router"
	"RESTAPI/internal/repository/sqlconnections"
	"RESTAPI/pkg/utils"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/net/http2"
)

// @title REST API Documentation
// @version 1.0
// @description This is a REST API for managing Teachers, Students and Executives
// @host localhost:443
// @BasePath /
// @schemes https
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT Bearer token
func main() {

	err := godotenv.Load("cmd/api/.env")
	if err != nil {
		log.Println("Error occured! ", err)
		return
	}
	log.Println("✅ .env file connected successfully")
	time.Sleep(time.Microsecond * 200)

	port := os.Getenv("API_PORT")
	cert := "cert.pem"
	key := "key.pem"

	_, err = sqlconnections.ConnectDb()
	if err != nil {
		log.Println("Error ------ ", err)
		return
	}
	log.Println("✅ Maria Database connected successfully")
	time.Sleep(time.Millisecond * 200)

	// Initialize Swagger documentation
	handler.InitSwagger()

	mux := router.Router()

	tlsConfige := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := mid.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortby", "SortOrder", "name", "age", "class", "first_name", "last_name", "subject", "id", "email", "page", "limit"},
	}

	rl := mid.NewLimiter(60, time.Second*30)
	securityMUX := utils.ApplayMiddlewares(mux, mid.Hpp(hppOptions), mid.Compress, mid.Exclude_Routes(mid.Security_middleware, "/swagger", "/swagger.json", "/execs/login", "/execs/logout", "/execs/forgotpassword", "/execs/resetpassword/reset/{resetcode}"), mid.Responce_time, rl.RL, mid.Exclude_Routes(mid.JWT_Middlerware, "/swagger", "/swagger.json", "/execs/login", "/execs/logout", "/execs/forgotpassword", "/execs/resetpassword/reset/{resetcode}"), mid.Sanitize)
	log.Println("✅ Security layers implemented  successfully")
	time.Sleep(time.Millisecond * 200)

	server := &http.Server{
		Addr:      port,
		TLSConfig: tlsConfige,
		Handler:   securityMUX,
	}

	http2.ConfigureServer(server, &http2.Server{})
	log.Println("✅ Server  implemented httpS successfully")
	time.Sleep(time.Millisecond * 200)

	log.Println("Server is running on port", port)
	server.ListenAndServeTLS(cert, key)
}
