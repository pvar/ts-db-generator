package tzdbio

import (
    "fmt"
    _ "github.com/mattn/go-sqlite3"
)

// GetZones retrieves available zones for specified timezone.
// The specified timezone is treated as a replica (link)
// which is first translated to the corresponding prototype timezone.
// The table of prototype timezones contains the name of the
// table with the corresponding zones.
func GetZones (timezone string) (zones []Zone, err error) {
    if !open {
        return nil, noConn
    }

    // get id of prototype timezone from replicas' table
    protoID, err := getReplicaPrototype (timezone)
    if err != nil {
        // cannot find prototype-id for specified replica
        return nil, err
    }

    // get all data for prototype timezone
    prototype, err := getPrototypeByID (protoID)
    if err != nil {
        // cannot find data for prototype with specified ID
        return nil, err
    }

    // check all available sub-tables with zones
    // start from the most recent -- the last one
    // stop when a reliable table is found
    tableOk := false
    for i := prototype.TabVer; i >= 0; i-- {
        zoneTable := fmt.Sprintf("%s%v", prototype.TabName, i)
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

    return
}

// getReplicaPrototype retrieves the prototype-ID for specified replica.
func getReplicaPrototype (replica string) (prototypeID int, err error) {
    columns := getReplicaCols();
    query := fmt.Sprintf("SELECT %s FROM %s WHERE %s=%s", columns[2], replicaTable, columns[1], replica)
    err = db.QueryRow(query).Scan(&prototypeID)

    if err != nil {
        return 0, err
    }

    return
}

// getPrototypeByID retrieves data for a prototype with specified ID.
func getPrototypeByID (prototypeID int) (prototype *Prototype, err error) {
    var name, zone, ztname string
    var id, ztver, offset int64

    columns := getPrototypeCols();
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%v", prototypeTable, columns[0], prototypeID)
    err = db.QueryRow(query).Scan(&id, &name, &zone, &offset, &ztname, &ztver)
    if err != nil {
        return nil, err
    }

    prototype = &Prototype{ID: id, Name: name, DZone: zone, DOffset: offset, TabName: ztname, TabVer: ztver}
    return
}

// getPrototypeByName retrieves ID for a named prototype.
func getPrototypeByName(prototypeName string) (*Prototype, error) {
    var name, dzone, ztname string
    var id, ztver, doffset int64

    columns := getPrototypeCols();
    query := fmt.Sprintf("SELECT * FROM %s WHERE %s=%v", prototypeTable, columns[1], prototypeName)
    err := db.QueryRow(query).Scan(&id, &name, &dzone, &doffset, &ztname, &ztver)
    if err != nil {
        return nil, err
    }

    return &Prototype{ID: id, Name: name, DZone: dzone, DOffset: doffset, TabName: ztname, TabVer: ztver}, nil
}

// getZones retrieves all zones from specified table.
func getZones (zoneTable string) (zones []Zone, err error) {
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

    return
}
