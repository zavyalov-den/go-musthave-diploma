package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
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

func (s *Storage) CreateOrder(ctx context.Context, num string, userID int) error {
	// language=sql
	query := `
		SELECT user_id from orders WHERE num = $1;
	`

	var storedUserID = 0

	err := s.db.QueryRow(ctx, query, num).Scan(&storedUserID)
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

	_, err = s.db.Exec(ctx, query, num, userID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Println("order already exists", pgErr.Code, pgErr.Message)
		} else {
			return err
		}
		//if errors.Is(err, pq) //
	}

	return nil
}

func (s *Storage) GetOrders(ctx context.Context, userID int) ([]*entities.Order, error) {
	// language=sql
	query := `
		SELECT num, status, accrual, uploaded_at from orders
		WHERE user_id = $1
	`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entities.Order

	for rows.Next() {
		order := &entities.Order{}
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}

	if len(result) == 0 {
		return nil, entities.ErrNoContent
	}

	return result, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order entities.Order) error {
	// language=sql
	query := `
		UPDATE orders SET status = $1, accrual = accrual + $2 WHERE num = $3
	`

	r, err := s.db.Exec(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}

	if r.RowsAffected() == 0 {
		return fmt.Errorf("order does not exist")
	}

	return nil

}
