package configs

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

var cfgFile string

var Cfg Config

var (
	defaults = map[string]interface{}{
		"port":  9000,
		"debug": false,
	}
	envPrefix   = ""
	configName  = "config"
	configType  = "yaml"
	configPaths = []string{
		".",
		fmt.Sprintf("%s/.oauth-showcase", getUserHomePath()),
	}
)

var allowedEnvVarKeys = []string{
	"port",
	"debug",
}

// Bootstrap reads in config file and ENV variables if set.
func Bootstrap() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		for _, p := range configPaths {
			viper.AddConfigPath(p)
		}
		viper.SetConfigType(configType)
		viper.SetConfigName(configName)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println(err)
		}
	}
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix(envPrefix)

	for _, key := range allowedEnvVarKeys {
		err := viper.BindEnv(key)
		if err != nil {
			fmt.Println(err)
		}
	}
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		log.Fatalf("could not decode config into struct: %v", err)
	}
}

func getUserHomePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func PrintConfig() {
	jsonConfig, err := json.MarshalIndent(&Cfg, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v\n", string(jsonConfig))
}

type Config struct {
	Port    int      `mapstructure:"port"`
	Debug   bool     `mapstructure:"debug"`
	Clients []Client `mapstructure:"clients"`
}

type Client struct {
	ID                          string   `mapstructure:"id"`
	Secret                      string   `mapstructure:"secret"`
	Domain                      string   `mapstructure:"domain"`
	RedirectUris                []string `mapstructure:"redirectUris"`
	AuthorizedGrantTypes        []string `mapstructure:"authorizedGrantTypes"`
	AccessTokenValiditySeconds  int      `mapstructure:"accessTokenValiditySeconds"`
	RefreshTokenValiditySeconds int      `mapstructure:"refreshTokenValiditySeconds"`
	AutoApproveScopes           []string `mapstructure:"autoApproveScopes"`
	Scopes                      []string `mapstructure:"scopes"`
	Users                       []User   `mapstructure:"users"`
}

type User struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
