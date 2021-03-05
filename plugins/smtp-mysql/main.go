package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/prometheus/alertmanager/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Export(val interface{}, params map[string]string) (interface{}, error) {
	c, ok := val.(*config.EmailConfig)
	if !ok {
		return nil, errors.New("*config.EmailConfig is required")
	}
	if db == nil {
		db0, err := sql.Open("mysql", params["uri"])
		if err != nil {
			return nil, err
		}
		db = db0
	}

	var user, pass string
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.QueryRowContext(ctx, params["query"]).Scan(&user, &pass)
	switch {
	case err == sql.ErrNoRows:
		return c, nil
	case err != nil:
		return nil, err
	default:
		c.AuthUsername = user
		c.AuthPassword = config.Secret(pass)
		return c, nil
	}
}

var _ config.Plugin = Export
