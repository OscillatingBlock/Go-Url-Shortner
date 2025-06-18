package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/go-sql-driver/mysql" // The Go SQL driver for MySQL (note the underscore!)
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect" // The Bun dialect for MySQL
)

type Url struct {
	bun.BaseModel `bun:"table:tiny,alias:u"`
	OriginalUrl   string `bun:"original_url,pk"`
	ShortenedUrl  string `bun:"shortened_url,unique,notnull"`
}

type DBManager struct {
	Db             *bun.DB
	ShortnerDomain string
}

func SetupDB(baseUrl string) DBManager {
	return DBManager{Db: setupDbInstance(), ShortnerDomain: baseUrl}
}

func setupDbInstance() *bun.DB {
	var db *bun.DB
	DSN := os.Getenv("DSN")
	slog.Info("DSN from .env", "dsn", DSN)
	if DSN == "" {
		slog.Error("Failed to get DSN from .env file")
	}
	sqldb, err := sql.Open("mysql", DSN)
	if err != nil {
		slog.Error("Failed to open MySQL connection: ", "error", err)
	}
	db = bun.NewDB(sqldb, mysqldialect.New())
	return db
}

func (manager *DBManager) CreateTables() error {
	ctx := context.Background()
	db := manager.Db
	_, err := db.NewCreateTable().Model((*Url)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to create url_parserer_db table: %w  \n", err)
	}
	return nil
}

func (manager *DBManager) get_url(shortCode string) (*Url, error) {
	var url Url
	ctx := context.Background()
	fullShortenedURL := manager.ShortnerDomain + shortCode
	err := manager.Db.NewSelect().Model(&url).Where("shortened_url = ?", fullShortenedURL).Scan(ctx)
	if err != nil {
		fmt.Printf("error while setting url %v \n", err)
		return nil, err
	}
	return &url, nil
}

func (manager *DBManager) set_url(original_url string) (*Url, error) {
	ctx := context.Background()
	// check in db if the url is already present
	var already_present_url Url
	err := manager.Db.NewSelect().Model(&already_present_url).Where("original_url= ?", original_url).Scan(ctx)
	if err == nil {
		return &already_present_url, err
	}
	// if not present then make new short url
	uniquePath := generateRandomString(SHORT_CODE_LEN)
	url := Url{OriginalUrl: original_url, ShortenedUrl: manager.ShortnerDomain + uniquePath}
	_, err2 := manager.Db.NewInsert().Model(&url).Exec(ctx)
	if err2 != nil {
		fmt.Printf("error while setting url %v \n", err)
		return nil, err
	}
	slog.Info("url set done")
	return &url, nil
}
