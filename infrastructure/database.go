package infrastructure

import (
	"database/sql"
	"github.com/kljensen/snowball"
	"os"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB() (*PostgresDB, error) {
	connStr := os.Getenv("POSTGRES_CONN")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &PostgresDB{db: db}, nil
}

func (p *PostgresDB) Close() error {
	return p.db.Close()
}

func (p *PostgresDB) GetResponse(intent string) (string, error) {
	var response string
	err := p.db.QueryRow("SELECT response FROM responses WHERE intent = $1", intent).Scan(&response)
	if err != nil {
		return "", err
	}
	return response, nil
}

func (p *PostgresDB) FindIntentByKeywords(query string) (string, string, error) {
	stemmedQuery, _ := snowball.Stem(query, "russian", true)

	var intent, response string
	err := p.db.QueryRow(
		`SELECT intent, response 
         FROM responses 
         ORDER BY similarity(intent, $1) DESC 
         LIMIT 1`, stemmedQuery,
	).Scan(&intent, &response)

	if err != nil {
		return "", "", err
	}
	return intent, response, nil
}

func (db *PostgresDB) GetAllIntents() ([]string, error) {
	rows, err := db.db.Query("SELECT DISTINCT intent FROM responses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var intents []string
	for rows.Next() {
		var intent string
		if err := rows.Scan(&intent); err != nil {
			return nil, err
		}
		intents = append(intents, intent)
	}
	return intents, nil
}
