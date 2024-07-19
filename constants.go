package main

const CreateTable string = `
	CREATE TABLE IF NOT EXISTS holdens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		name TEXT,
    	age INTEGER,
		height REAL
  );
`
const getQuery string = "SELECT * FROM holdens WHERE id = ?"

const getAllQuery string = "SELECT * FROM holdens"
