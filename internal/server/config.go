package server

import (
	"flag"
	"fmt"
	"github.com/Nikolay961996/goferma/internal/utils"
	"github.com/caarlos0/env/v6"
	"os"
)

type Config struct {
	runAddress           string
	databaseURI          string
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
	flag.StringVar(&c.runAddress, "a", "localhost:8080", "Run address")
	flag.StringVar(&c.databaseURI, "d", "", "Database address")
	flag.StringVar(&c.accrualSystemAddress, "r", "", "Accrual system address")
	flag.StringVar(&c.secretKey, "k", "MY_SUPER_SECRET_KEY", "Secret key for signing")

	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		utils.Log.Info(fmt.Sprintf("  -%s: %v (default: %v)\n", f.Name, f.Value, f.DefValue))
	})

	if flag.NArg() > 0 {
		utils.Log.Fatal("To many args!")
	}
	utils.Log.Info("Address from flag: ", c.runAddress)
	utils.Log.Info("Database from flag: ", c.databaseURI)
}

func (c *Config) parseEnv() {
	var envConfig struct {
		RunAddress           string `env:"RUN_ADDRESS"`
		DatabaseUri          string `env:"DATABASE_URI"`
		AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
		SecretKey            string `env:"SIGNING_SECRET_KEY"`
	}

	err := env.Parse(&envConfig)
	if err != nil {
		utils.Log.Fatal(err)
	}

	utils.Log.Info("Address from env: ", envConfig.RunAddress)
	utils.Log.Info("Database from env: ", envConfig.DatabaseUri)

	test := os.Getenv("DATABASE_URI")
	utils.Log.Info("test: ", test)

	if envConfig.RunAddress != "" {
		c.runAddress = envConfig.RunAddress
	}
	if envConfig.DatabaseUri != "" {
		c.databaseURI = envConfig.DatabaseUri
	}
	if envConfig.AccrualSystemAddress != "" {
		c.accrualSystemAddress = envConfig.AccrualSystemAddress
	}
	if envConfig.SecretKey != "" {
		c.secretKey = envConfig.SecretKey
	}
}

func (c *Config) check() {
	if c.runAddress == "" {
		utils.Log.Fatal("Run address is required")
	}
	if c.databaseURI == "" {
		utils.Log.Fatal("Database address is required")
	}
	if c.accrualSystemAddress == "" {
		utils.Log.Fatal("Accrual system address is required")
	}
}
