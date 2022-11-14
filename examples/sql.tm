#!/usr/bin/env tamarin

conn := sql.connect("postgres://postgres:mysecretpassword@localhost:5432/postgres")

result := sql.query(conn, "select * from users where age > $1", 40)

assert(result.is_ok(), "SQL query failed")

rows := result.unwrap()
rows.each(func(row) { print("row:", row) })
print("returned", len(rows), "matching rows")

conn | sql.query("select COUNT(*), NOW() from users")
