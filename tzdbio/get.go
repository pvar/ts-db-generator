package tzdbio

import (
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// GetZones retrieves available zones for specified timezone.
// The specified timezone is treated as a replica (link)
// which is first translated to the corresponding original TZ.
// The table of original timezones contains the name of the
// table with the corresponding zones.
func GetZones (timezone string) (zones []Zone, err error) {
    if !dbOpen {
        return nil, noDB
    }

    // get id of original timezone from replicas' table
    protoID, err := getReplicaOriginal (timezone)
    if err != nil {
        // cannot find original TZ for specified replica
        return nil, err
    }

    // get all data for original timezone
    original, err := getOriginalByID (protoID)
    if err != nil {
        // cannot find data for original TZ
        return nil, err
    }

    // check all available sub-tables with zones
    // start from the most recent -- the last one
    // stop when a reliable table is found
    tableOk := false
    for i := original.TabVer; i > 0; i-- {
        zoneTable := fmt.Sprintf("%s%v", original.TabName, i)
        zones, err = getZones (zoneTable)
        if err != nil {
            // zone table unreliable
            continue
        }
        tableOk = true
        break
    }

    if !tableOk {
        return nil, fmt.Errorf("tzdbio: cannot find reliable table with zones")
    }

    return zones, nil
}

// getReplicaOriginal retrieves the original-ID for specified replica.
func getReplicaOriginal (replicaTZ string) (originalID int, err error) {
    columns := getReplicaCols();
    query := fmt.Sprintf("SELECT %s FROM %s WHERE %s=%q", columns[2], replicaTable, columns[1], replicaTZ)
    err = db.QueryRow(query).Scan(&originalID)

    if err != nil {
        return 0, err
    }

    return originalID, nil
}

// getOriginalByID retrieves data for an origial TZ with specified ID.
func getOriginalByID (originalID int) (*Original, error) {
    var name, zone, ztname, tzdver string
    var id, ztver, offset int64

    columns := getOriginalCols();
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%v", originalTable, columns[0], originalID)
    err := db.QueryRow(query).Scan(&id, &name, &zone, &offset, &ztname, &ztver, &tzdver)
    if err != nil {
        return nil, err
    }

    return &Original{ID: id, Name: name, DZone: zone, DOffset: offset, TabName: ztname, TabVer: ztver, TZDVer: tzdver}, nil
}

// getOriginalByName retrieves ID for a named origial TZ.
func getOriginalByName(originalTZ string) (*Original, error) {
    var name, dzone, ztname, tzdatver string
    var id, tzver, doffset int64

    columns := getOriginalCols();
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%q", originalTable, columns[1], originalTZ)
    err := db.QueryRow(query).Scan(&id, &name, &dzone, &doffset, &ztname, &tzver, &tzdatver)
    if err != nil {
        return nil, err
    }

    return &Original{ID: id, Name: name, DZone: dzone, DOffset: doffset, TabName: ztname, TabVer: tzver, TZDVer: tzdatver}, nil
}

// getZones retrieves all zones from specified table.
func getZones (zoneTable string) (zones []Zone, err error) {
    if !tableExists(zoneTable) {
        return nil, err
    }

    query := fmt.Sprintf("SELECT * FROM %s", zoneTable)
    rows, err := db.Query(query)
    defer  rows.Close()
    if err != nil {
        return nil, err
    }

    var id, start, end, offset int64
    var name string
    var isDST bool

    zones = make([]Zone, 0, 5)
    for rows.Next() {
        err = rows.Scan(&id, &name, &start, &end, &offset, &isDST)
        if err != nil {
            return nil, err
        }
        newZone := Zone{ID: id, Name: name, Start: start, End: end, Offset: offset, IsDST: isDST}
        zones = append(zones, newZone)
    }

    return zones, nil
}
