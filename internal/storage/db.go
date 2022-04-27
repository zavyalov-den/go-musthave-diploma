package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
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

func (s *Storage) Register(ctx context.Context, cred *entities.Credentials) (int, error) {
	// language=sql
	var userID int
	query := `
		INSERT INTO users(login, password) VALUES ($1, $2)
		returning id;
	`
	err := s.db.QueryRow(ctx, query, cred.Login, cred.Password).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *Storage) GetUser(ctx context.Context, name string) (*entities.Credentials, error) {
	var user entities.Credentials
	// language=sql
	query := `
		SELECT id, login, password FROM users
		WHERE login = $1
	`
	err := s.db.QueryRow(ctx, query, name).Scan(&user.UserID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Storage) CreateOrder(ctx context.Context, num int, userID int) error {
	// language=sql
	query := `
		SELECT user_id from orders WHERE num = $1;
	`

	var storedUserID = 0

	err := s.db.QueryRow(ctx, query, string(num)).Scan(&storedUserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if storedUserID != 0 && storedUserID != userID {
		return entities.ErrUserConflict
	}

	// language=sql
	query = `
		insert into orders (num, user_id) 
		values ($1, $2)
	`

	_, err = s.db.Exec(ctx, query, string(num), userID)
	if err != nil {
		return err
	}

	return nil
}
