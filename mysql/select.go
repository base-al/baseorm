package baseorm

import (
	"dbconn"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func all[T any]() ([]T, error) {
	var results []T
	t := reflect.TypeOf((*T)(nil)).Elem()

	tableName := strings.ToLower(t.Name()) + "s"
	var fields []string
	for i := 0; i < t.NumField(); i++ {
		field := strings.ToLower(t.Field(i).Name)
		fields = append(fields, field)
	}

	query := fmt.Sprintf("SELECT %s FROM `%s`;", strings.Join(fields, ", "), tableName)
	log.Printf("Executing query: %s\n", query)

	rows, err := dbconn.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		itemPtr := reflect.New(t).Interface()
		scanArgs := make([]interface{}, t.NumField())
		for i := range scanArgs {
			scanArgs[i] = reflect.ValueOf(itemPtr).Elem().Field(i).Addr().Interface()
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		results = append(results, reflect.ValueOf(itemPtr).Elem().Interface().(T))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
