package configs

import (
	"fmt"
	"log"

	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)


type conf struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBHost string `mapstructure:"DB_HOST"`
	DBPort string `mapstructure:"DB_PORT"`
	DBUser string `mapstructure:"DB_USER"`
	DBName string `mapstructure:"DB_NAME"`
	DBPass string `mapstructure:"DB_PASS"`
	WebServerPort string `mapstructure:"WEB_SERVER_PORT"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	JWTExpireIn int  `mapstructure:"JWT_EXPIRE_IN"`
	TokenAuth *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
    cfg := &conf{}

    viper.AddConfigPath(path)
    viper.SetConfigFile(".env")
    viper.AutomaticEnv()

    err := viper.ReadInConfig()
    if err != nil {
        log.Printf("Error while reading config file %s", err)
        return nil, err
    }

    err = viper.Unmarshal(cfg)
    if err != nil {
        log.Printf("Error while unmarshalling config %s", err)
        return nil, err
    }

    cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

    return cfg, nil
}


func (c *conf) GetDBDSN() string {
    return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC",
        c.DBHost, c.DBPort, c.DBUser, c.DBName, c.DBPass)
}