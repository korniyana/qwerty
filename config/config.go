package config

import "github.com/spf13/viper"

const (
	configName = "qwerty"
	configExt  = "yaml"
)

type Config struct {
	Logger *Logger `mapstructure:"logger"`
	Server *Server `mapstructure:"server"`
}

type Logger struct {
	File            string `mapstructure:"file"`
	Format          string `mapstructure:"format"`
	Level           string `mapstructure:"level"`
	Log2Engine      int    `mapstructure:"log2engine"`
	TimeStampFormat string `mapstructure:"time_stamp_format"`
}

type Server struct {
	PidFile     string `mapstructure:"pid_file"`
	Name        string `mapstructure:"name"`
	BindAddress string `mapstructure:"bind_address"`
}

// Load reads configuration from file or environment variables.
func Load(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configExt)

	//viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	return cfg, viper.Unmarshal(cfg)
}
