package server

import (
	"flag"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	runAddress           string
	databaseUri          string
	accrualSystemAddress string
	secretKey            string
}

func NewConfig() *Config {
	c := &Config{}
	c.parseFlags()
	c.parseEnv()
	c.check()

	return c
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.runAddress, "a", "", "Run address")
	flag.StringVar(&c.databaseUri, "d", "", "Database address")
	flag.StringVar(&c.accrualSystemAddress, "r", "", "Accrual system address")
	flag.StringVar(&c.secretKey, "k", "MY_SUPER_SECRET_KEY", "Secret key for signing")

	flag.Parse()

	if flag.NArg() > 0 {
		utils.Log.Fatal("To many args!")
	}
	utils.Log.Info("Address from flag: ", c.runAddress)
	utils.Log.Info("Database from flag: ", c.databaseUri)
}

func (c *Config) parseEnv() {
	var envConfig struct {
		runAddress           string `env:"RUN_ADDRESS"`
		databaseUri          string `env:"DATABASE_URI"`
		accrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
		secretKey            string `env:"SIGNING_SECRET_KEY"`
	}

	err := env.Parse(&envConfig)
	if err != nil {
		utils.Log.Fatal(err)
	}

	utils.Log.Info("Address from env: ", envConfig.runAddress)
	utils.Log.Info("Database from env: ", envConfig.databaseUri)

	test := os.Getenv("DATABASE_URI")
	utils.Log.Info("test: ", test)

	if envConfig.runAddress != "" {
		c.runAddress = envConfig.runAddress
	}
	if envConfig.databaseUri != "" {
		c.databaseUri = envConfig.databaseUri
	}
	if envConfig.accrualSystemAddress != "" {
		c.accrualSystemAddress = envConfig.accrualSystemAddress
	}
	if envConfig.secretKey != "" {
		c.secretKey = envConfig.secretKey
	}
}

func (c *Config) check() {
	if c.runAddress == "" {
		utils.Log.Fatal("Run address is required")
	}
	if c.databaseUri == "" {
		utils.Log.Fatal("Database address is required")
	}
	if c.accrualSystemAddress == "" {
		utils.Log.Fatal("Accrual system address is required")
	}
}
