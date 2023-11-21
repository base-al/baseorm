package baseorm

import (
	"dbconn"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func update[T any](item T) error {
	t := reflect.TypeOf(item)

	tableName := strings.ToLower(t.Name()) + "s"
	var fields []string
	var values []string
	for i := 0; i < t.NumField(); i++ {
		field := strings.ToLower(t.Field(i).Name)
		fields = append(fields, field)
		values = append(values, "?")
	}

	query := fmt.Sprintf("UPDATE `%s` SET %s WHERE id = ?;", tableName, strings.Join(fields, "= ?, "))
	log.Printf("Executing query: %s\n", query)

	stmt, err := dbconn.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var args []interface{}
	for i := 1; i < t.NumField(); i++ {
		args = append(args, reflect.ValueOf(item).Field(i).Interface())
	}
	args = append(args, reflect.ValueOf(item).Field(0).Interface())

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}