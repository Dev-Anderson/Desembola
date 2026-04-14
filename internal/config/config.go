package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type AppConfig struct {
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSslMode           string
	CSVFileName         string
	CleanDBBeforeInsert bool
	APIPort             string
}

func LoadConfig(envFile string) *AppConfig {
	loadEnv(envFile)

	csvFile := os.Getenv("CSV_FILE")
	if csvFile == "" {
		csvFile = "teste.csv"
	}

	return &AppConfig{
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBSslMode:           os.Getenv("DB_SSLMODE"),
		CSVFileName:         csvFile,
		CleanDBBeforeInsert: strings.ToLower(os.Getenv("CLEAN_DB_BEFORE_INSERT")) == "true",
		APIPort:             os.Getenv("API_PORT"),
	}
}

// loadEnv lê um arquivo de texto simples no formato CHAVE=VALOR e define
// as variáveis de ambiente equivalentes na execução atual.
func loadEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Println("Aviso: arquivo .env não encontrado, prosseguindo com variáveis do sistema.")
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
		}
	}
}
