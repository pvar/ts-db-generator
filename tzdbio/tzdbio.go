package tzdbio

import (
    "errors"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB
    open bool
    noConn = errors.New("tzdbio: no connection to any db file")
)

func init () {
    open = false
}

func openDB (filename string) error {
    dbObj, err := sql.Open("sqlite", filename)

    if err != nil {
        open = false
    } else {
        open = true
        db = dbObj
    }

    return err
}

func closeDB () error {
    if !open {
        return noConn
    }

    db.Close()
    return nil
}
