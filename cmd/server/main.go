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
	"gorm.io/driver/sqlite"
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
	//SQLite
	db, err := gorm.Open(sqlite.Open(config.DBPath), &gorm.Config{})
//	PostgreSQL
// 	db, err := gorm.Open(postgres.Open(config.GetDBDSN()), &gorm.Config{})
	
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{},&entity.Project{},&entity.Tenant{},&entity.Chain{},&entity.Table{},&entity.Rule{})

	//Debug config values JWTExpireIn
	log.Println("Config JWTExpireIN: " + strconv.FormatUint(uint64(config.JWTExpireIn), 10))
	userHandler := handlers.NewUserHandler(database.NewUser(db))
	tenantHandler := handlers.NewTenantHandler(database.NewTenantDB(db))
	projectHandler := handlers.NewProjectHandler(database.NewProjectDB(db), database.NewTenantDB(db))
	tableHandler := handlers.NewTableHandler(database.NewTableDB(db))
	chainHandler := handlers.NewChainHandler(database.NewChainDB(db), database.NewProjectDB(db), database.NewTableDB(db))
	ruleHandler := handlers.NewRuleHandler(database.NewRuleDB(db), database.NewChainDB(db), database.NewServiceDB(db), database.NewNetworkObjectDB(db))
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
		r.Get("/filter", tenantHandler.GetTenantsWithFilters)
        r.Get("/{id}", tenantHandler.GetTenantByID)
		r.Get("/name/{name}", tenantHandler.GetTenantByName)
        r.Get("/", tenantHandler.GetTenants)
        r.Put("/{id}", tenantHandler.UpdateTenant)
        r.Delete("/{id}", tenantHandler.DeleteTenant)
    })
	r.Route("/projects", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", projectHandler.CreateProject)
		r.Get("/{id}", projectHandler.GetProject)
		r.Get("/filter", projectHandler.GetProjectsWithFilters)
		r.Put("/{id}", projectHandler.UpdateProject)
		r.Delete("/{id}", projectHandler.DeleteProject)
		r.Get("/", projectHandler.GetProjects)
	})
	r.Route("/tables", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", tableHandler.CreateTable)
		r.Get("/{id}", tableHandler.GetTable)
		r.Get("/filter", tableHandler.GetTablesWithFilters)
		r.Put("/{id}", tableHandler.UpdateTable)
		r.Delete("/{id}", tableHandler.DeleteTable)
		r.Get("/", tableHandler.GetTables)
	})
	r.Route("/chains", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", chainHandler.CreateChain)
		r.Get("/{id}", chainHandler.GetChain)
		r.Get("/filter", chainHandler.GetChainsWithFilters)
		r.Put("/{id}", chainHandler.UpdateChain)
		r.Delete("/{id}", chainHandler.DeleteChain)
		r.Get("/", chainHandler.GetChains)
	})
	r.Route("/rules", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", ruleHandler.CreateRule)
		// r.Get("/{id}", ruleHandler.GetRule)
		// r.Get("/filter", ruleHandler.GetRulesWithFilters)
		// r.Put("/{id}", ruleHandler.UpdateRule)
		// r.Delete("/{id}", ruleHandler.DeleteRule)
		// r.Get("/", ruleHandler.GetRules)
	})


		
	r.Post("/users/generate_token", userHandler.GetJWT)

	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8000/docs/swagger/doc.json")))

	http.ListenAndServe(":"+config.WebServerPort, r)
}
