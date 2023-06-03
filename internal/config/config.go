package config

import "github.com/spf13/viper"

// Функция инициализации конфига (всех токенов)
func InitConfig() error {

	// Путь и имя файла
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()

}
