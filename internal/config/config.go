package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBName               string        `mapstructure:"DB_NAME"`
	AdminUserName        string        `mapstructure:"ADMIN_USER_NAME"`
	AdminPassword        string        `mapstructure:"ADMIN_PASSWORD"`
	Environment          string        `mapstructure:"APP_ENV"`
	InDocker             string        `mapstructure:"IN_DOCKER"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func (c *Config) DBSource() string {
	return fmt.Sprintf("data/%s", c.DBName)
}
func (c *Config) DBSourceURL() string {
	return fmt.Sprintf("sqlite3://%s", c.DBSource())
}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AllowEmptyEnv(true)
	v.SetDefault("DB_NAME", "app.db")
	v.SetDefault("ADMIN_USER_NAME", "admin")
	v.SetDefault("ADMIN_PASSWORD", "admin123")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("IN_DOCKER", "false")
	v.SetDefault("TOKEN_SYMMETRIC_KEY", "9y$B&E)H@McQfTjWnZr4u7x!A%D*G-Ka")
	v.SetDefault("ACCESS_TOKEN_DURATION", time.Minute*5)
	v.SetDefault("REFRESH_TOKEN_DURATION", time.Hour*24*30)
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
