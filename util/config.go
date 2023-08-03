package util

import (
	"github.com/spf13/viper"

	"time"
)

type Config struct {
	// DB
	DBDriver      string `mapstructure:"DB_DRIVER" env:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE" env:"DB_SOURCE"`
	DBUrl         string `mapstructure:"DB_URL" env:"DB_URL"`
	RedisUrl      string `mapstructure:"REDIS_URL" env:"REDIS_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD" env:"REDIS_PASSWORD"`
	RedisUsername string `mapstructure:"REDIS_USERNAME" env:"REDIS_USERNAME"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS" env:"SERVER_ADDRESS"`

	// PubNub Key
	PubNubPublishKey   string `mapstructure:"PUBNUB_PUBLISH_KEY" env:"PUBNUB_PUBLISH_KEY"`
	PubNubSubscribeKey string `mapstructure:"PUBNUB_SUBSCRIBE_KEY" env:"PUBNUB_SUBSCRIBE_KEY"`
	PubNubUserId       string `mapstructure:"PUBNUB_USER_ID" env:"PUBNUB_USER_ID"`

	// HTTP address
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS" env:"GRPC_SERVER_ADDRESS"`
	GRPCClientAddress string `mapstructure:"GRPC_CLIENT_ADDRESS" env:"GRPC_CLIENT_ADDRESS"`
	SocketAddress     string `mapstructure:"SOCKET_ADDRESS" env:"SOCKET_ADDRESS"`

	// JWT
	TokenSymmetricKey      string        `mapstructure:"TOKEN_SYMMETRIC_KEY" env:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration    time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" env:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY" env:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY" env:"ACCESS_TOKEN_PRIVATE_KEY"`
	RefreshTokenDuration   time.Duration `mapstructure:"REFRESH_TOKEN_DURATION" env:"REFRESH_TOKEN_DURATION"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY" env:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY" env:"REFRESH_TOKEN_PRIVATE_KEY"  `
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
