package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//go:embed schema.sql
var schemaSQL string

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

func DefaultConfig() Config {
	return Config{
		User:     "root",
		Password: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		DBName:   "zang_ye",
	}
}

func New(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库连接: %w", err)
	}

	database.SetMaxOpenConns(10)
	database.SetMaxIdleConns(5)
	database.SetConnMaxLifetime(5 * time.Minute)

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("数据库连接失败: %w", err)
	}

	if _, err := database.Exec(schemaSQL); err != nil {
		database.Close()
		return nil, fmt.Errorf("建表失败: %w", err)
	}

	return database, nil
}