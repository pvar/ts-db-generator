package tzdbio

import (
    "fmt"
)

// AddFullOriginal adds data to an existing entry in table of origial TZs.
// This function is mainly used during initial setup, after having parsed
// and processed the respective timezone file.
func AddFullOriginal (origTZ *Original) error {
    if !dbOpen {
        return noDB
    }

    fields := getOriginalCols()
    query := fmt.Sprintf("UPDATE %q SET %q=?, %q=?, %q=?, %q=?, %q=? WHERE %q=%q",
                originalTable, fields[2], fields[3], fields[4], fields[5],
                fields[6], fields[1], origTZ.Name)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(origTZ.DZone, origTZ.DOffset, origTZ.TabName, origTZ.TabVer, origTZ.TZDVer)
    return err
}

// AddOriginal adds the name of a new entry in table of original TZs.
// The rest of the data remain uninitialized. This function is used
// during initial setup, to populate table with available origials.
func AddOriginal (originalTZ string) (id int64, err error) {
    if !dbOpen {
        return -1, noDB
    }

    fields := getOriginalCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s, %s) VALUES(?, ?, ?, ?, ?, ?)",
                originalTable, fields[1], fields[2], fields[3], fields[4], fields[5], fields[6])

    stmt, err := db.Prepare(query)
    if err != nil {
        return -1, err
    }
    defer stmt.Close()

    res, err := stmt.Exec(originalTZ, "", -1, "", -1, "")
    if err != nil {
        return -1, err
    }

    id, err = res.LastInsertId()
    if err != nil {
        // this is highly impropable,
        // since the DB statement was
        // executed without any errors...
        return -1, err
    }

    return id, nil
}

// AddReplicas adds a new list of entries in the preplicas' table.
// Each group of replicas contains the name of the original as an
// extra entry. This function is mainly used during initial setup,
// to populate table with replicas.
func AddReplicas (replicaTZs []string, originalTZ string) error {
    if !dbOpen {
        return noDB
    }

    id, err := needOriginalID(originalTZ)
    if err != nil {
        return err
    }

    fields := getReplicaCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES(?, ?)",
                replicaTable, fields[1], fields[2])

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // add each replica with the ID of the specified origial TZ
    for _, replicaTZ := range replicaTZs {
        _, err := stmt.Exec(replicaTZ, id)
        if err != nil {
            return err
        }
    }

    return nil
}

// AddZones adds a new sub-table of zones, to the specified origial timezone.
func AddZones (originalTZ string, zones []Zone) error {
    if !dbOpen {
        return noDB
    }

    origial, err := getOriginalByName(originalTZ)
    if err != nil {
        return err
    }

    newZonesTable := fmt.Sprintf("%s%v", origial.TabName, origial.TabVer + 1)
    createTable (newZonesTable)

    fields := getZoneCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s) VALUES(?, ?, ?, ?, ?)",
                newZonesTable, fields[1], fields[2], fields[3], fields[4], fields[5])

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    var dst int
    for _, zone := range zones {
        if zone.IsDST {
            dst = 1
        } else {
            dst = 0
        }
        _, err := stmt.Exec(zone.Name, zone.Start, zone.End, zone.Offset, dst)
        if err != nil {
            return err
        }
    }

    return nil
}

// UpdateOriginal updates default zone and offset of an origial timezone.
func UpdateOriginal (origTZ *Original) error {
    if !dbOpen {
        return noDB
    }

    fields := getOriginalCols()
    query := fmt.Sprintf("UPDATE %s SET %s=? %s=? WHERE %s=%q",
                originalTable, fields[2], fields[3], fields[1], origTZ.Name)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(origTZ.DZone, origTZ.DOffset)

    return err
}

// UpdateReplica updates the origial timezone linked to the specified replica.
func UpdateReplica (replicaTZ, originalTZ string) error {
    if !dbOpen {
        return noDB
    }

    id, err := needOriginalID(originalTZ)
    if err != nil {
        return err
    }

    fields := getReplicaCols()
    query := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=%q",
                replicaTable, fields[2], fields[1], replicaTZ)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    return err
}

// needOriginalID retrieves ID for named origial TZ or creates it.
func needOriginalID (originalTZ string) (id int64, err error) {
    origial, err := getOriginalByName(originalTZ)
    if err != nil {
        // Could not get ID for specified original timezone.
        // Attempt to add it and get ID of new entry.
        id, err = AddOriginal (originalTZ)
        if err != nil {
            return -1, err
        }
        return id, nil
    }

    return origial.ID, nil
}
