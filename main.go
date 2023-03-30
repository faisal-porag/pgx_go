package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
)

// define a struct to represent the table columns
type MyTableRow struct {
	Column1 string
	Column2 string
}

func main() {
	// create a connection to the database
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:admin@localhost:5432/mydb")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	err1 := conn.Ping(context.Background())
	if err1 != nil {
		log.Println(err1)
	}

	rows := [][]interface{}{
		{"John", "Smith"},
		{"Jane", "Doe"},
	}

	copyCount, err := conn.CopyFrom(
		context.Background(),
		pgx.Identifier{"mytable"},
		[]string{"column1", "column2"},
		pgx.CopyFromRows(rows),
	)
	
	log.Println(copyCount)

}

func listTasks(conn *pgx.Conn) error {
	rows, _ := conn.Query(context.Background(), "select * from tasks")

	for rows.Next() {
		var id int32
		var description string
		err := rows.Scan(&id, &description)
		if err != nil {
			return err
		}
		fmt.Printf("%d. %s\n", id, description)
	}

	return rows.Err()
}

func addTask(description string, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "insert into tasks(description) values($1)", description)
	return err
}

func updateTask(itemNum int32, description string, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "update tasks set description=$1 where id=$2", description, itemNum)
	return err
}

func removeTask(itemNum int32, conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), "delete from tasks where id=$1", itemNum)
	return err
}
