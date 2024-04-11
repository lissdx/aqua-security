package providers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/lissdx/aqua-security/internal/drivers"
	asLogger "github.com/lissdx/aqua-security/internal/pkg/logger"
	"github.com/spf13/viper"
)

func NewStore(logger asLogger.TOLogger, config *viper.Viper) drivers.Store {
	dbFullURL := fmt.Sprintf("%s/%s?sslmode=%s", config.GetString("DB_CONNECTION_URL"), config.GetString("DB_NAME"), config.GetString("DB_SSL_MODE"))
	conn, err := sql.Open(config.GetString("DB_DRIVER"), dbFullURL)
	if err != nil {
		logger.Fatal("cannot connect to db:", err)
	}

	return drivers.NewStore(conn, logger)
}
