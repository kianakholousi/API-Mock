package config

import "github.com/spf13/viper"

type Config struct {
	Database *Database
	Server   *Server
	Security *Security
}

type Database struct {
	Driver   string
	Host     string
	Port     int
	DB       string
	Username string
	Password string
	Charset  string
}

type Server struct {
	Host string
	Port int
}

type Security struct {
	SecretKey           string
	EncryptionAlgorithm string
}

type Params struct {
	FilePath string
	FileName string
	FileType string
}

func Init(param Params) (*Config, error) {
	viper.SetConfigType(param.FileType)
	viper.SetConfigFile(param.FileName)
	viper.AddConfigPath(param.FilePath)

	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	database := &Database{
		Driver:   viper.GetString("database.driver"),
		Host:     viper.GetString("database.host"),
		Port:     viper.GetInt("database.port"),
		DB:       viper.GetString("database.db"),
		Username: viper.GetString("database.username"),
		Password: viper.GetString("database.password"),
		Charset:  viper.GetString("database.charset"),
	}

	server := &Server{
		Host: viper.GetString("server.host"),
		Port: viper.GetInt("server.port"),
	}

	security := &Security{
		SecretKey:           viper.GetString("security.secret_key"),
		EncryptionAlgorithm: viper.GetString("security.encryption_algorithm"),
	}

	return &Config{
		Database: database,
		Server:   server,
		Security: security,
	}, nil
}
