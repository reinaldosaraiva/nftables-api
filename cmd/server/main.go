package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/reinaldosaraiva/nftables-api/configs"
	_ "github.com/reinaldosaraiva/nftables-api/docs"
	"github.com/reinaldosaraiva/nftables-api/internal/entity"
	"github.com/reinaldosaraiva/nftables-api/internal/infra/database"
	"github.com/reinaldosaraiva/nftables-api/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title	NFTABLES API Go
// @version	0.1
// @description	This is a Go API that provides a Restful interface for managing nftables, a powerful and flexible firewall framework in the Linux kernel.
// @termsOfService	http://swagger.io/terms
// @contact.name	Reinaldo Saraiva
// @contact.email	reinaldo.saraiva@gmail.com
// license.name		Bearware
// host				localhost:8000
// @BasePath	/
// @securityDefinitions.apikey	ApiKeyAuth
// @in	header
// @name	Authorization
func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	
	db, err := gorm.Open(postgres.Open(config.GetDBDSN()), &gorm.Config{})
	
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{},&entity.Project{},&entity.Tenant{},&entity.Chain{},&entity.Table{},&entity.Rule{})

	//Debug config values JWTExpireIn
	log.Println("Config JWTExpireIN: " + strconv.FormatUint(uint64(config.JWTExpireIn), 10))
	userHandler := handlers.NewUserHandler(database.NewUser(db))
	tenantHandler := handlers.NewTenantHandler(database.NewTenantDB(db))
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JWTExpireIn))
	r.Post("/users", userHandler.CreateUser)
	r.Get("/users", userHandler.GetUserByEmail)
	r.Route("/tenants", func(r chi.Router) { // Tenant routes
        r.Use(jwtauth.Verifier(config.TokenAuth))
        r.Use(jwtauth.Authenticator)
        r.Post("/", tenantHandler.CreateTenant)
        r.Get("/{id}", tenantHandler.GetTenant)
        r.Get("/", tenantHandler.GetTenants)
        r.Put("/{id}", tenantHandler.UpdateTenant)
        r.Delete("/{id}", tenantHandler.DeleteTenant)
    })
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/swagger/doc.json")))

	http.ListenAndServe(":"+config.WebServerPort, r)
}
