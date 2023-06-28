#!/usr/bin/env risor

func connect(user="postgres", pass="", host="localhost", port=5432, db="postgres") {
    return pgx.connect('postgres://{user}:{pass}@{host}:{port}/{db}')
}

conn := connect("postgres", "mysecretpassword")
if conn.is_err() {
    print(conn.err_msg())
    exit(1)
}

conn.query("select * from users where age > $1", 18).
    each(func(row) {
        print("row:", row)
    })
