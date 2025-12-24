package config

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppPort              string        `mapstructure:"APP_PORT"`
	DBHost               string        `mapstructure:"DB_HOST"`
	DBPort               string        `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPassword           string        `mapstructure:"DB_PASSWORD"`
	DBName               string        `mapstructure:"DB_NAME"`
	Environment          string        `mapstructure:"APP_ENV"`
	InDocker             string        `mapstructure:"IN_DOCKER"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func (c *Config) DBSource() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

}
func (c *Config) DBSourceURL() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", c.DBUser, url.QueryEscape(c.DBPassword), c.DBHost, c.DBPort, c.DBName)

}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	v.SetDefault("APP_PORT", "8123")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", "5432")
	v.SetDefault("DB_USER", "hanzo")
	v.SetDefault("DB_PASSWORD", "WVO574bJJAtr")
	v.SetDefault("DB_NAME", "bizbundl")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("IN_DOCKER", "false")
	v.SetDefault("TOKEN_SYMMETRIC_KEY", "9y$B&E)H@McQfTjWnZr4u7x!A%D*G-Ka")
	v.SetDefault("ACCESS_TOKEN_DURATION", time.Minute*5)
	v.SetDefault("REFRESH_TOKEN_DURATION", time.Hour*24*30)

	// Bind environment variables
	bindEnvs(v, Config{})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	return &cfg
}
func bindEnvs(v *viper.Viper, iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		vf := ifv.Field(i)
		tf := ift.Field(i)
		tv := tf.Tag.Get("mapstructure")
		if tv == "" {
			continue
		}
		switch vf.Kind() {
		case reflect.Struct:
			bindEnvs(v, vf.Interface(), append(parts, tv)...)
		default:
			_ = v.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
