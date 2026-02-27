package pkg

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"main/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	*sql.DB
}

func NewPostgresDatabase(logger Logger, env config.Env) PostgresDB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.PGHost,
		env.PGPort,
		env.PGUser,
		env.PGPass,
		env.PGName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logger.Fatal("Failed to open PostgreSQL connection: ", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("Failed to ping PostgreSQL: ", err)
	}

	logger.Info("Successfully connected to PostgreSQL")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal("migration driver error: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+env.MigrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		logger.Fatal("migration init error: ", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal("migration failed: ", err)
	}

	logger.Info("Successfully applied database migrations")

	return PostgresDB{db}
}

func (db PostgresDB) Close() error {
	return db.DB.Close()
}

func (db PostgresDB) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
