package main

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

func InitDB(driverName string, databaseName string, tableCreate string) *sql.DB {
	val, err := sql.Open(driverName, databaseName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = val.Exec(tableCreate)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

func GetOne[T any](db *sql.DB, query string, params ...interface{}) (*T, error) {
	row := db.QueryRow(query, params...)

	// Get the type of the generic parameter T
	elemType := reflect.TypeOf((*T)(nil)).Elem()
	elem := reflect.New(elemType).Elem()

	// Create a slice to hold the column values
	columnValues := make([]interface{}, elemType.NumField())
	columnPointers := make([]interface{}, elemType.NumField())
	for i := 0; i < elemType.NumField(); i++ {
		columnPointers[i] = &columnValues[i]
	}

	// Scan the row into the column pointers
	if err := row.Scan(columnPointers...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Map the column values to the struct fields
	for i := 0; i < elemType.NumField(); i++ {
		field := elem.Field(i)
		val := reflect.ValueOf(columnValues[i])
		if val.Type().ConvertibleTo(field.Type()) {
			field.Set(val.Convert(field.Type()))
		}
	}

	// Return the populated struct
	result := elem.Addr().Interface().(*T)
	return result, nil
}

func GetList[T any](db *sql.DB, query string) ([]T, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	// Get column names from the query result
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Prepare a slice to hold the results
	var results []T

	for rows.Next() {
		// Create a new instance of type T
		elemType := reflect.TypeOf((*T)(nil)).Elem()
		elem := reflect.New(elemType).Elem()

		// Create a slice of interfaces to hold the column values
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// Scan the row into the column pointers
		if err = rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		// Map the column values to the struct fields
		for i := 0; i < elemType.NumField(); i++ {
			field := elem.Field(i)
			val := reflect.ValueOf(columnValues[i])
			if val.Type().ConvertibleTo(field.Type()) {
				field.Set(val.Convert(field.Type()))
			}
		}

		// Append the result to the results slice
		results = append(results, elem.Interface().(T))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
