package database

import "database/sql"

type Database struct {
	sql *sql.DB
}

func (db *Database) setUserTimeStamp(userId string, initialTimeStamp string) {
	
}