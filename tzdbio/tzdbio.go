package tzdbio

import (
    "fmt"
    "strings"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

var (
    db *sql.DB
    dbOpen bool
    noDB = fmt.Errorf("tzdbio: no connection to db")
)

func init () {
    dbOpen = false
}

func Open (filename string) error {
    dbObj, err := sql.Open("sqlite3", filename)

    if err != nil {
        dbOpen = false
        return err
    }

    dbOpen = true
    db = dbObj

    if !tableExists(originalTable) {
        createTable(getOriginalSchema())
    }

    if !tableExists(replicaTable) {
        createTable(getReplicaSchema())
    }

    return nil
}

func Close () error {
    if !dbOpen {
        return noDB
    }

    db.Close()
    return nil
}

func tableExists(tableName string) bool {
    var tempname string
    query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='{%s}';", tableName)
    row := db.QueryRow(query)
    err := row.Scan(&tempname)
    if err != nil {
        return false
    }

    return true
}

func createTable(query string) error {
    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    _, err = stmt.Exec()
    if err != nil {
        return err
    }

    return nil
}

func createZoneTable(tableName string) error {
    query := getZoneSchema(tableName)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    _, err = stmt.Exec()
    if err != nil {
        return err
    }

    return nil
}

func makeTabName (prototype string, version int) (tableName string, err error) {
    if len(prototype) == 0 {
        return "", fmt.Errorf("Original TZ name is empty!")
    }

    if version < 0 {
        return "", fmt.Errorf("Invalid version number!")
    }

    r := strings.NewReplacer("/", "_", "\\", "_", "+", "_P_", "-", "_M_")
    base := r.Replace(prototype)
    base = strings.ToLower(base)
    return fmt.Sprintf("%s%d", base, version), nil
}
