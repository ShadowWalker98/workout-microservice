package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

func (app *application) connectDB() *sql.DB {
	conn, err := sql.Open("postgres", app.config.db.dsn)
	if err != nil {
		app.logger.Println("Error while connecting to database: ", err)
		return nil
	}

	conn.SetMaxOpenConns(app.config.db.maxOpenConns)
	conn.SetMaxIdleConns(app.config.db.maxIdleConns)
	duration, err := time.ParseDuration(app.config.db.maxIdleTime)
	if err != nil {
		return nil
	}

	conn.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = conn.PingContext(ctx)
	if err != nil {
		app.logger.Println(err)
		return nil
	}
	return conn
}
