package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func (s *Storage) initDB(ctx context.Context) {
	//language=sql
	var queries = []string{`
	CREATE TABLE IF NOT EXISTS users (
	    id serial primary key,
	    username text not null,
	    password text not null
	);
`,
		`
	CREATE TABLE IF NOT EXISTS orders (
	    id serial primary key,
	    num text unique
	    
	);
`,
		`
	CREATE TABLE IF NOT EXISTS user_orders (
	    user_id int references users(id),
	    order_id int references orders(id)
	)
`}
	for _, q := range queries {
		_, err := s.db.Exec(ctx, q)
		if err != nil {
			panic(err)
		}
	}
}
