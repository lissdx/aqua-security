package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

var dbNotExistErrorFormat = "database \"%s\" does not exist"

func CreateDBIfNotExist(cfg MigrationConfig) error {
	db, err := sql.Open(cfg.DbDriver, cfg.FullURL())
	if err != nil {
		_ = db.Close()
		return fmt.Errorf("connecting to db: %s", err)
	}

	//defer db.Close()

	exists, err := dbExists(db, cfg.DbName)
	if err != nil {
		return err
	}

	// close this DB. we don't need it anymore
	if err = db.Close(); err != nil {
		return err
	}

	if !exists {
		log.Printf("creating the new db: %s ", cfg.DbName)
		if err = createDB(cfg); err != nil {
			return fmt.Errorf("creating db: %s", err)
		}
		log.Printf("... the new db is created: %s ", cfg.DbName)
	} else {
		log.Printf("... the db %s is existed. do nothing  ", cfg.DbName)
		return db.Close()
	}

	return nil
}

func dbExists(db *sql.DB, name string) (bool, error) {
	//row := db.QueryRow(fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_database WHERE datname = '%s')", name))
	//
	//var exists bool
	//err := row.Scan(&exists)
	//if err != nil {
	//	return false, err
	//}
	err := db.Ping()
	if err != nil {
		if strings.HasSuffix(err.Error(), fmt.Sprintf(dbNotExistErrorFormat, name)) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func createDB(cfg MigrationConfig) error {
	// point to default postgres DB
	db, err := sql.Open(cfg.DbDriver, cfg.URLWithDefaultDb())
	defer db.Close()

	if err != nil {
		return fmt.Errorf("connecting to db: %s", err)
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DbName)); err != nil {
		return err
	}

	return nil
}
