package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/ThomasNguyenGitHub/go/log"
	"github.com/ThomasNguyenGitHub/go/storage/local"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	gl "gorm.io/gorm/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseAuth struct {
	Host     string
	Port     string
	UserName string
	Password string
	DBName   string
	SSLMode  string
}

func Connect(opts ...gorm.Option) (*gorm.DB, error) {

	var (
		host         = local.Getenv("PG_DB_HOST")
		port         = local.Getenv("PG_DB_PORT")
		dbName       = local.Getenv("PG_DB_NAME")
		user         = local.Getenv("PG_DB_USER")
		password     = local.Getenv("PG_DB_PASS")
		pgSqlConnStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=prefer",
			host, port, user, password, dbName)
		options = []gorm.Option{
			&gorm.Config{
				Logger: gl.Default.LogMode(gl.Silent),
			},
		}
	)
	options = append(options, opts...)
	return gorm.Open(postgres.Open(pgSqlConnStr), options...)
}

func Connection() (db *sql.DB, err error) {
	connStr, err := getDBConfigs()
	if err != nil {
		log.Printf("Cannot get db confuration: %s", err)
		return nil, err
	}
	connector, err := pq.NewConnector(connStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db = sql.OpenDB(connector)
	log.Printf("We are connected to the %s database", "postgres")
	defer db.Close()
	return db, err
}

func getDBConfigs() (string, error) {
	var dbInfo = DatabaseAuth{}
	dbInfo = DatabaseAuth{
		Host:     local.Getenv("PG_DB_HOST"),
		Port:     local.Getenv("PG_DB_PORT"),
		UserName: local.Getenv("PG_DB_USER"),
		Password: local.Getenv("PG_DB_PASS"),
		DBName:   local.Getenv("PG_DB_NAME"),
		SSLMode:  local.Getenv("PG_DB_SSL_MODE"),
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbInfo.UserName, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.DBName, dbInfo.SSLMode)
	log.Print(connStr)
	return connStr, nil

}
