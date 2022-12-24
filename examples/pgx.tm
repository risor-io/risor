#!/usr/bin/env tamarin

conn := pgx.connect("postgres://postgres:mysecretpassword@localhost:5432/postgres")
result := conn.query("select * from users where age > $1", 18)
result.each(func(row) { print("row:", row) })

conn.query("select COUNT(*), NOW() from users")
