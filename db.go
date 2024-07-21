package main

import (
	"database/sql"
	"errors"
	"log"
	"reflect"
)

type RowScanner interface {
	Scan(dest ...interface{}) error
}

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
	return getOne[T](db.QueryRow(query, params...))
}

func GetList[T any](db *sql.DB, query string, params ...interface{}) ([]*T, error) {
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	var results []*T
	for rows.Next() {
		res, rowErr := getOne[T](rows)
		if rowErr != nil {
			return nil, rowErr
		}
		results = append(results, res)
	}
	return results, nil
}

func getOne[T any](scanner RowScanner) (*T, error) {
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
	if err := scanner.Scan(columnPointers...); err != nil {
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
