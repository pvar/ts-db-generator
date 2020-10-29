package tzdbio

import (
    "fmt"
//    "errors"
)

// AddPrototype adds *only* the name of a new entry in prototypes' table.
// The rest of the data remain uninitialized. This function is mainly used
// during initial setup, to populate table with available prototypes.
func AddPrototype (prototypeName string) (id int64, err error) {
    if !open {
        return -1, noConn
    }

    fields := getPrototypeCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES(?, ?)", prototypeTable, fields[1], fields[5])

    stmt, err := db.Prepare(query)
    if err != nil {
        return -1, err
    }
    defer stmt.Close()

    res, err := stmt.Exec(prototypeName, -1)
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

    return
}

// AddReplicas adds a new list of entries in the preplicas' table.
// Each group of replicas contains the name of the original as an
// extra entry. This function is mainly used during initial setup,
// to populate table with replicas.
func AddReplicas (replicas []string, prototypeName string) error {
    if !open {
        return noConn
    }

    id, err := needPrototypeID(prototypeName)
    if err != nil {
        return err
    }

    fields := getReplicaCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES(?, ?)", replicaTable, fields[1], fields[2])

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    // add each replica with the ID of the specified prototype
    for _, replica := range replicas {
        _, err := stmt.Exec(replica, id)
        if err != nil {
            return err
        }
    }

    return nil
}

// AddZones adds a new sub-table of zones, to the specified prototype timezone.
func AddZones (prototypeName string, zones []Zone) error {
    if !open {
        return noConn
    }

    prototype, err := getPrototypeByName(prototypeName)
    if err != nil {
        return err
    }

    newZonesTable := fmt.Sprintf("%s%v", prototype.TabName, prototype.TabVer + 1)
    createTable (newZonesTable)

    fields := getZoneCols()
    query := fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s, %s) VALUES(?, ?, ?, ?, ?)", newZonesTable, fields[1], fields[2], fields[3], fields[4], fields[5])

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

// UpdatePrototype updates default zone and offset of a prototype timezone.
func UpdatePrototype (prototype Prototype) error {
    if !open {
        return noConn
    }

    fields := getPrototypeCols()
    query := fmt.Sprintf("UPDATE %s SET %s=? %s=? WHERE %s=%s",
            prototypeTable, fields[2], fields[3], fields[1], prototype.Name)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(prototype.DZone, prototype.DOffset)

    return err
}

// UpdateReplica updates the prototype timezone linked to the specified replica.
func UpdateReplica (replicaName, prototypeName string) error {
    if !open {
        return noConn
    }

    id, err := needPrototypeID(prototypeName)
    if err != nil {
        return err
    }

    fields := getReplicaCols()
    query := fmt.Sprintf("UPDATE %s SET %s=? WHERE %s=%s", replicaTable, fields[2], fields[1], replicaName)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    return err
}

// needPrototypeID retrieves ID for named prototype or creates it.
func needPrototypeID (prototypeName string) (id int64, err error) {
    prototype, err := getPrototypeByName(prototypeName)
    if err != nil {
        // Could not get ID for specified prototype.
        // Attempt to add it and get ID of new entry.
        id, err = AddPrototype (prototypeName)
        if err != nil {
            return -1, err
        }
        return id, nil
    }

    return prototype.ID, nil
}
