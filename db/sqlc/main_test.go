package db

//mainly testing whether it can connect to postgresql
import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db", err)
	}
	defer conn.Close()
	// conn.SetMaxIdleConns(3)
	// conn.SetMaxOpenConns(3)
	testQueries = New(conn) //new function defined in db.go file
	os.Exit(m.Run())        //it will tell golng test files are successful to run or fail
}
