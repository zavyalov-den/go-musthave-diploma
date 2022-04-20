package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/zavyalov-den/go-musthave-diploma/internal/config"
	"github.com/zavyalov-den/go-musthave-diploma/internal/entities"
	"log"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage() *Storage {
	cfg, err := pgxpool.ParseConfig(config.GetConfig().DatabaseURI)
	if err != nil {
		log.Fatal("failed to parse db config: ", err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), cfg)
	if err != nil {
		log.Fatal("failed to connect to db: ", err)
	}

	return &Storage{db: db}

}

func (s *Storage) Register(ctx context.Context, cred entities.Credentials) error {
	// language=sql
	query := `
		INSERT INTO users(username, password) VALUES ($1, $2);
	`
	res, err := s.db.Exec(ctx, query, cred.Username, cred.Password)
	if err != nil {
		return err
	}
	rows := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("user wasn't created")
	}

	return nil
}

func (s *Storage) GetUser(ctx context.Context, name string) (*entities.Credentials, error) {
	var user entities.Credentials
	// language=sql
	query := `
		SELECT username, password FROM users
		WHERE username = $1
	`
	err := s.db.QueryRow(ctx, query, name).Scan(&user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil

}
