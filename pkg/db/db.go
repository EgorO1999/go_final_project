package db

import (
    "database/sql"
    "os"
    _ "modernc.org/sqlite"
)

var DB *sql.DB

const sqlScript = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(255) NOT NULL DEFAULT '',
    comment TEXT DEFAULT '',
    repeat VARCHAR(128) DEFAULT ''
);

CREATE INDEX idx_scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {
    _, err := os.Stat(dbFile)
    install := false
    if err != nil {
        install = true
    }

    db, err := sql.Open("sqlite", dbFile)
    if err != nil {
        return err
    }

    if install {
        if _, err := db.Exec(sqlScript); err != nil {
            db.Close()
            return err
        }
    }

    DB = db
    return nil
}