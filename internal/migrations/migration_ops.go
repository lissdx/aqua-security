package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/viper"
	"log"
)

const defaultPostgresDbName = "postgres"

func RunDBMigration(cfg MigrationConfig) bool {

	migration, err := migrate.New(cfg.MigrationURL, cfg.FullURL())
	if err != nil {
		log.Printf("creating migrate instance: %s\n", err)
		return false
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("running migration: %s\n", err)
		return false
	}

	return true
}

type MigrationConfig struct {
	DbDriver     string
	DbName       string
	DbURL        string
	DbSSLMode    string
	MigrationURL string
}

func (mc *MigrationConfig) FullURL() string {
	return fmt.Sprintf("%s/%s?sslmode=%s", mc.DbURL, mc.DbName, mc.DbSSLMode)
}

func (mc *MigrationConfig) URLWithDefaultDb() string {
	return fmt.Sprintf("%s/%s?sslmode=%s", mc.DbURL, defaultPostgresDbName, mc.DbSSLMode)
}

func NewMigrationConfig(envCfg *viper.Viper) MigrationConfig {
	return MigrationConfig{
		DbDriver:     envCfg.GetString("DB_DRIVER"),
		DbName:       envCfg.GetString("DB_NAME"),
		DbURL:        envCfg.GetString("DB_CONNECTION_URL"),
		DbSSLMode:    envCfg.GetString("DB_SSL_MODE"),
		MigrationURL: envCfg.GetString("MIGRATION_URL"),
	}
}

//func NewAuthMigrationConfig(envCfg *viper.Viper) MigrationConfig {
//	return MigrationConfig{
//		DbDriver:     envCfg.GetString("DB_DRIVER"),
//		DbName:       envCfg.GetString("AUTH_DB_NAME"),
//		DbURL:        envCfg.GetString("DB_CONNECTION_URL"),
//		DbSSLMode:    envCfg.GetString("DB_SSL_MODE"),
//		MigrationURL: envCfg.GetString("AUTH_MIGRATION_URL"),
//	}
//}
