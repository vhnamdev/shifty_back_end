package configs

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// ------- App Configuration--------
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`

	// --------Postgres Database Configuration-----------
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// --------Mongo Database Configuration-----------
	MongoURI    string `mapstructure:"MONGO_URI"`
	MongoDBName string `mapstructure:"MONGO_DB_NAME"`

	// --------Redis Database Configuration-----------
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// --------Qdrant Database Configuration-----------
	QdrantHost     string `mapstructure:"QDRANT_HOST"`
	QdrantPort     string `mapstructure:"QDRANT_PORT"`
	QdrantGrpcPort string `mapstructure:"QDRANT_GRPC_PORT"`

	// --------Security(JWT & OAuth) Configuration-----------
	JWTAccessSecret    string `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	JWTRefreshSecret   string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
	AcessTokenExpiry   string `mapstructure:"ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry string `mapstructure:"REFRESH_TOKEN_EXPIRY"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleSecret       string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleCallBackUrl  string `mapstructure:"GOOGLE_CALLBACK_URL"`
}

func LoadConfig() (config *Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Can not find file .env")
	}
	err = viper.Unmarshal(&config)
	return
}
