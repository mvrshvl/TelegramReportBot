package database

import (
	"TelegramBot/config"
	"context"
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"io/fs"
	"net/http"
	"sync"
	"time"
)

const (
	dsn = "%s:%s@tcp(%s)/%s?parseTime=true"
)

//go:embed migrations/*.sql
var migrationsPath embed.FS

type Database struct {
	mux        sync.RWMutex
	connection *sqlx.DB

	cfg *config.Database
}

func New(cfg *config.Config) *Database {
	return &Database{cfg: &cfg.Database}
}

func (db *Database) Connect(ctx context.Context) error {
	err := db.connect(ctx, db.cfg, dsn)
	if err != nil {
		return err
	}

	return db.migrate(db.cfg, migrate.Up)
}

func (db *Database) GetConnection() *sqlx.DB {
	return db.connection
}

func (db *Database) migrate(cfg *config.Database, direction migrate.MigrationDirection) error {
	migrationsDirectory, err := fs.Sub(migrationsPath, "migrations")
	if err != nil {
		return err
	}

	migrations := &migrate.HttpFileSystemMigrationSource{FileSystem: http.FS(migrationsDirectory)}

	_, err = migrate.Exec(db.connection.DB, cfg.Driver, migrations, direction)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) connect(ctx context.Context, cfg *config.Database, dsn string) (err error) {
	db.connection, err = sqlx.ConnectContext(ctx, cfg.Driver, fmt.Sprintf(dsn, cfg.User, cfg.Password, cfg.Address, cfg.Name))
	if err != nil {
		return fmt.Errorf("failed to connect db: %w", err)
	}

	return nil
}

func (db *Database) Read(ctx context.Context, city string) ([]*Message, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	query := `SELECT * FROM messages
				WHERE city = ?
				ORDER BY timestamp DESC`

	rows, err := db.connection.QueryContext(ctx, query, city)
	if err != nil {
		db.reconnect()
		return nil, fmt.Errorf("can't get accounts: %w", err)
	}

	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		var message Message

		err := rows.Scan(
			&message.ID,
			&message.City,
			&message.Text,
			&message.Media,
			&message.Place,
			&message.From,
			&message.Timestamp,
		)
		if err != nil {
			db.reconnect()
			return nil, err
		}

		messages = append(messages, &message)
	}

	return messages, nil
}

func (db *Database) Insert(msg *Message) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	query := `INSERT INTO MESSAGES (city, text, media, place, from, timestamp) VALUES (?, ?, ?, ?, ?, ?)`

	_, err := db.connection.Exec(query, msg.City, msg.Text, msg.Media, msg.Place, msg.From, msg.Timestamp)
	if err != nil {
		db.reconnect()
		return fmt.Errorf("can't insert message: %w", err)
	}

	return nil
}

func (db *Database) reconnect() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db.mux.Lock()
	defer db.mux.Unlock()

	ticker := time.NewTimer(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		err := db.connect(ctx, db.cfg, dsn)
		if err != nil {
			fmt.Println(err)

			continue
		}

		return
	}
}
