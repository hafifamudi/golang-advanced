package db

import (
	"database/sql"
	"fmt"
	"log"
)

type database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) *database {
	return &database{db: db}
}

func (d database) Migration() {
	migrationSQL := `
	create table if not exists goqite (
	  id text primary key default ('m_' || lower(hex(randomblob(16)))),
	  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
	  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
	  queue text not null,
	  body blob not null,
	  timeout text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
	  received integer not null default 0
	) strict;

	create trigger if not exists goqite_updated_timestamp after update on goqite begin
	  update goqite set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
	end;

	create index if not exist goqite_queue_created_idx on goqite (queue, created);
	`

	_, err := d.db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("Error executing schema migration: %v", err)
	}
	fmt.Println("Schema migration executed successfully")
}
