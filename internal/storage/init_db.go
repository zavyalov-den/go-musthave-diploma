package storage

import "context"

func (s *Storage) InitDB() {
	//language=sql
	var queries = []string{`
	CREATE TABLE IF NOT EXISTS users (
	    id serial primary key,
	    username text not null unique,
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
		_, err := s.db.Exec(context.Background(), q)
		if err != nil {
			panic(err)
		}
	}
}
