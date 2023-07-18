package config

import (
	"chatto/internal/util/strutil"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

const (
	accessTokenDuration  = time.Minute * 60
	refreshTokenDuration = time.Hour * 24 * 90
)

type AppConfig struct {
	Ip              string `mapstructure:"LISTEN_IP"`
	Port            string `mapstructure:"LISTEN_PORT"`
	Address         string
	UserDatabaseURI string `mapstructure:"USER_DB_URI"`

	ChatDatabaseURI      string `mapstructure:"CHAT_DB_URI"`
	ChatDatabaseUsername string `mapstructure:"CHAT_DB_USERNAME"`
	ChatDatabasePassword string `mapstructure:"CHAT_DB_PASSWORD"`

	JWTSecretKey    string `mapstructure:"JWT_SECRET_KEY"`
	JWTSigningType  string `mapstructure:"JWT_SIGNING_TYPE"`
	JWTSecretKeyURI string `mapstructure:"JWT_SECRET_KEY_URI"`

	// Duration is in minutes
	AccessTokenDuration  uint64 `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration uint64 `mapstructure:"REFRESH_TOKEN_DURATION"`

	JWTKeyFunc jwt.Keyfunc
}

func LoadConfig(path string) (AppConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.SetDefault("LISTEN_IP", "localhost")
	viper.SetDefault("LISTEN_PORT", 9999)
	viper.SetDefault("ACCESS_TOKEN_DURATION", accessTokenDuration)
	viper.SetDefault("REFRESH_TOKEN_DURATION", refreshTokenDuration)

	if err := viper.ReadInConfig(); err != nil {
		return AppConfig{}, err
	}

	log.Println("Config File Used:" + viper.ConfigFileUsed())

	var conf AppConfig
	err := viper.Unmarshal(&conf)
	conf.Address = conf.Ip + ":" + conf.Port

	if strutil.IsEmpty(conf.UserDatabaseURI) {
		return conf, errors.New("database uri should not be absent, set DB_URI on env")
	}
	if strutil.IsEmpty(conf.JWTSecretKey) && strutil.IsEmpty(conf.JWTSecretKeyURI) {
		return conf, errors.New("JWT secret key should not be absent, set JWT_SECRET_KEY on env")
	}
	if strutil.IsEmpty(conf.JWTSigningType) {
		return conf, errors.New("JWT Signing type should not be absent, set JWT_SIGNING_TYPE on env")
	}

	// Set the function to get the secret key either by the config or response
	if len(conf.JWTSecretKeyURI) == 0 {
		conf.JWTKeyFunc = func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.JWTSecretKey), nil
		}
	} else {
		conf.JWTKeyFunc = func(token *jwt.Token) (interface{}, error) {
			resp, err := http.Get(conf.JWTSecretKeyURI)
			defer resp.Body.Close()
			if err != nil {
				return nil, err
			}

			var data []byte
			_, err = resp.Body.Read(data)
			if err != nil {
				return nil, err
			}

			return data, nil
		}
	}

	return conf, err
}
