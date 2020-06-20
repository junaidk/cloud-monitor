package config

import (
	"github.com/spf13/viper"
	"log"
)

var cfg *Config
var ConfPath string

func (conf *Config) vp() error {

	v := viper.New()
	//v.AddConfigPath(".")      // path to look for the config file in
	//v.AddConfigPath(ConfPath) // path to look for the config file in

	v.SetConfigFile(ConfPath)
	v.AutomaticEnv()

	err := v.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Println(err)
		return err
	}

	err = v.Unmarshal(conf)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type Config struct {
	Azure AzureConfig `mapstructure:"azure"`
	DO    DOConfig    `mapstructure:"do"`
	GCP   GCPConfig   `mapstructure:"gcp"`
	AWS   AWSConfig   `mapstructure:"aws"`
	IBM   IBMConfig   `mapstructure:"ibm"`
}

type IBMConfig struct {
	APIKey string `mapstructure:"api_key"`
}

type AzureConfig struct {
	ID           string `mapstructure:"client_id"`
	Key          string `mapstructure:"client_secret"`
	Tenant       string `mapstructure:"tenant_id"`
	Subscription string `mapstructure:"subscription_id"`
}
type DOConfig struct {
	AccessToken string `mapstructure:"access_token"`
}
type GCPConfig struct {
	CredentialFilePath string `mapstructure:"credential_file"`
}
type AWSConfig struct {
	AccessKey    string `mapstructure:"access_key"`
	AccessSecret string `mapstructure:"access_secret"`
}

func (conf *Config) ParseConfig() error {
	//log.Println("parsing configurations")
	err := conf.vp()
	return err
}

func NewConfig() (*Config, error) {

	if cfg != nil {
		return cfg, nil
	} else {
		cfg = &Config{}
		err := cfg.ParseConfig()
		return cfg, err
	}

}
