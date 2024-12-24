package utils

import (
	"github.com/spf13/viper"
)

// Config stores all configuration of the application
// The values are read by viper from a config file or environment variables
type Config struct {
	AMQPURL       string `mapstructure:"AMQP_URL"`
	QueueName     string `mapstructure:"QUEUE_NAME"`
	FromEmail     string `mapstructure:"FROM_EMAIL"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      string `mapstructure:"SMTP_PORT"`
}

// LoadConfig reads configuration file or environment variables.
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
