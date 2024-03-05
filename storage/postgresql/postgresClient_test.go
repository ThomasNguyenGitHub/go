package postgresql

import (
	"database/sql"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"testing"
)

func TestConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	connStr, err := getDBConfigs()
	if err != nil {
		t.Fatal(err)
	}
	connector, err := pq.NewConnector(connStr)
	if err != nil {
		t.Fatal(err)
	}
	db := sql.OpenDB(connector)
	defer db.Close()
	if err := db.PingContext(ctx); err != nil {
		t.Fatal("expected Ping to succeed")
	}
	defer cancel()
	txn, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	txn.Rollback()

}
