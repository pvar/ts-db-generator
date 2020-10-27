package tzdbio

import (
    "fmt"
    "errors"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB
    open bool
    noConn  = errors.New("tzdbio: no connection to db")
    noValid = errors.New("tzdbio: database seems empty")
)

func init () {
    open = false
}

func Open (filename string) error {
    dbObj, err := sql.Open("sqlite3", filename)

    if err != nil {
        open = false
        return err
    }

    open = true
    db = dbObj

    err = checkTable(prototypeTable)
    if err != nil {
        createTable(getPrototypeSchema)
    }

    err := checkTable(replicaTable)
    if err != nil {
        createTable(getReplicaSchema)
    }

    return nil
}

func Close () error {
    if !open {
        return noConn
    }

    db.Close()
    return nil
}

func checkTable(tableName string) error {
    var tempname string
    query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='{%s}';", tableName)
    row := db.QueryRow(query)
    err := row.Scan(&tempname)
    return err
}
