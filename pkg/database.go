package pkg

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/muhammadtaufan/go-sensor-collector/config"
)

func InitDatabase(cfg *config.Config) (*sqlx.DB, error) {
	dsnConn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DATABASE_USERNAME, cfg.DATABASE_PASSWORD, cfg.DATABASE_HOST, cfg.DATABASE_NAME)
	db, err := sqlx.Open("mysql", dsnConn)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL database: %w", err)
	}

	return db, nil
}
