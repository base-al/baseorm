// Description: Delete data from database
// Version: 0.1
// Requirements:
// 	- dbconn package
// 	- reflect package
// 	- strings package
// 	- fmt package
// 	- log package
// 	- database/sql package
// 	- github.com/go-sql-driver/mysql package
// 	- T any
// 	- item T
// Input: item T

package baseorm

import (
	"dbconn"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/postgresql"
)

func delete(item interface{}) error {
	t := reflect.TypeOf(item)

	tableName := strings.ToLower(t.Name()) + "s"
	var fields []string
	var values []string
	for i := 0; i < t.NumField(); i++ {
		field := strings.ToLower(t.Field(i).Name)
		fields = append(fields, field)
		values = append(values, "?")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", tableName)
	log.Printf("Executing query: %s\n", query)

	stmt, err := dbconn.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var args []interface{}
	args = append(args, reflect.ValueOf(item).Field(0).Interface())

	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
