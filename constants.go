package main

const CreateTable string = `
	CREATE TABLE IF NOT EXISTS holdens (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
  		name TEXT,
    	age INTEGER,
		height REAL
  );
`
