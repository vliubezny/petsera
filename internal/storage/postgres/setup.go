package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratep "github.com/golang-migrate/migrate/v4/database/postgres"

	// init file source
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host, Port                 string
	User, Password             string
	DBName                     string
	MaxOpenConns, MaxIdleConns int
	EnableMigration            bool
	MigrationsDir              string
}

// MustSetupDB opens DB connection and runs migrations.
func MustSetupDB(dsn string, maxOpenConns, maxIdleConns int, enableMigration bool, migrations string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create postgres connection")
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	if err := db.PingContext(context.Background()); err != nil {
		logrus.WithError(err).Fatal("failed to ping postgres")
	}

	if !enableMigration {
		return db
	}

	driver, err := migratep.WithInstance(db, &migratep.Config{})
	if err != nil {
		logrus.WithError(err).Fatal("failed to create database migrate driver")
	}

	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrations), "", driver)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create migrator")
	}

	checkVersion(migrator)

	switch err := migrator.Up(); err {
	case nil:
		logrus.Info("database was migrated")
	case migrate.ErrNoChange:
		logrus.Info("database is up-to-date")
	default:
		logrus.WithError(err).Fatal("failed to migrate db")
	}

	checkVersion(migrator)

	return db
}

func checkVersion(m *migrate.Migrate) {
	switch v, d, err := m.Version(); err {
	case nil:
		logrus.Infof("database version %d with dirty state %t", v, d)
	case migrate.ErrNilVersion:
		logrus.Info("database version: nil")
	default:
		logrus.WithError(err).Fatal("failed to get version")
	}
}

func newDSN(config Config) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBName,
	)
}

func MustSetupStorage(config Config) PG {
	dsn := newDSN(config)
	db := MustSetupDB(
		dsn,
		config.MaxOpenConns,
		config.MaxIdleConns,
		config.EnableMigration,
		config.MigrationsDir,
	)

	return New(db)
}
