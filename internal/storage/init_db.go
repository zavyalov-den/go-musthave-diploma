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
	CREATE TABLE IF NOT EXISTS balance (
	    id serial primary key,
	    user_id int unique references users(id),
	    current float default 0,
	    withdrawn float default 0
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
	    order_num text,
	    sum float default 0,
	    processed_at timestamptz default now()
	)
`}
	for _, q := range queries {
		_, err := s.db.Exec(context.Background(), q)
		if err != nil {
			panic(err)
		}
	}
}
