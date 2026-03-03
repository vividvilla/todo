package db

import "database/sql"

func SetConfig(key, value string) error {
	_, err := DB.Exec(
		`INSERT INTO config (key, value) VALUES (?, ?)
		 ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	return err
}

func GetConfig(key string) (string, error) {
	var value string
	err := DB.QueryRow(`SELECT value FROM config WHERE key = ?`, key).Scan(&value)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}
