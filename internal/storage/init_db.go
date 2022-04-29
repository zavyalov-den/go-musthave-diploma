package storage

import "context"

func (s *Storage) InitDB() {
	//language=sql
	var queries = []string{`
	CREATE TABLE IF NOT EXISTS users (
	    id serial primary key,
	    login text not null unique,
	    password text not null
	);
`,
		`
	CREATE TABLE IF NOT EXISTS orders (
	    id serial primary key,
	    num text unique,
	    user_id int references users(id),
		status text default 'NEW',
	    accrual float default 0,
	    uploaded_at timestamptz default now()
	);
`,
		`
	CREATE TABLE IF NOT EXISTS withdrawals (
	    id serial primary key,
	    user_id int references users(id),
	    order_id int references orders(id),
	    amount int
	)
`}
	for _, q := range queries {
		_, err := s.db.Exec(context.Background(), q)
		if err != nil {
			panic(err)
		}
	}
}
