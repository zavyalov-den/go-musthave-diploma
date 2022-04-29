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

	// language=sql
	query = `
		INSERT INTO balance(user_id) VALUES ($1)
		returning id;
	`
	_, err = s.db.Exec(ctx, query, userID)
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

func (s *Storage) GetUserBalance(ctx context.Context, userID int) (*entities.Balance, error) {
	var balance entities.Balance
	// language=sql
	query := `
		SELECT current, withdrawn FROM balance
		WHERE user_id = $1;
	`
	err := s.db.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return &balance, nil
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
	} else if storedUserID != 0 && storedUserID == userID {
		return entities.ErrEntryExists

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

func (s *Storage) UpdateOrder(ctx context.Context, order entities.AccrualOrder) error {
	// language=sql
	query := `
		UPDATE orders SET status = $1, accrual = accrual + $2 WHERE num = $3
	`

	_, err := s.db.Exec(ctx, query, order.Status, order.Accrual, order.Order)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateUserBalance(ctx context.Context, userID int, accrual float32) error {
	// language=sql
	query := `
		UPDATE balance SET current = current + $1 WHERE user_id = $2
	`

	_, err := s.db.Exec(ctx, query, accrual, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) Withdraw(ctx context.Context, userID int, withdrawal entities.Withdrawal) error {
	// language=sql
	query := `
		INSERT INTO withdrawals(user_id, order_num, sum)
		VALUES ($1, $2, $3)
		returning id;
	`

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, query, userID, withdrawal.Order, withdrawal.Sum)
	if err != nil {
		return err
	}

	// language=sql
	query = `
		UPDATE balance SET withdrawn = withdrawn - $1 WHERE user_id = $2
	`

	_, err = tx.Exec(ctx, query, withdrawal.Sum, userID)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
