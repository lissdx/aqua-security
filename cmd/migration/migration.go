package main

import (
	"github.com/lissdx/aqua-security/internal/migrations"
	asConfig "github.com/lissdx/aqua-security/internal/pkg/config"
	"log"
	"os"
)

const ver = "20240410"

func main() {
	envCfg := asConfig.NewConfig()

	migrationCfg := migrations.NewMigrationConfig(envCfg)
	log.Printf("run migration ver: %s\n", ver)
	log.Printf("migration cnf: %+v\n", migrationCfg)

	// because CREATE DB is potentially very dangerous operation
	// I've added a protection level
	// it means we must explicitly allow it
	if envCfg.GetBool("RUN_CREATE_NEW_DB") {
		log.Println("Try to create DB if not exists")
		if err := migrations.CreateDBIfNotExist(migrationCfg); err != nil {
			log.Printf("cannot create the new DB. migrationConfig: %+v\nerror: %s", migrationCfg, err.Error())
			os.Exit(1)
		}
	} else {
		log.Println("create DB process is not allowed. Skip it...")
	}

	log.Println("Run the migration process")
	success := migrations.RunDBMigration(migrationCfg)
	if !success {
		log.Println("Migration Failure")
		os.Exit(1)
	}

	log.Println("Migration Success")
	os.Exit(0)
}
