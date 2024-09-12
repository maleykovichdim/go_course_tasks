module maleykovich.db/db/pgsql

go 1.22.3

replace maleykovich.db/db => ../

require (
	github.com/jackc/pgx/v5 v5.6.0
	maleykovich.db/db v0.0.0-00010101000000-000000000000
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/text v0.14.0 // indirect
)
