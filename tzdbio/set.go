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
    query := fmt.Sprintf("INSERT INTO prototypes (%s, %s) VALUES(?, ?)", fields[1], fields[5])

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

// UpdatePrototype updates default zone and offset of a prototype timezone.
func UpdatePrototype (prototype Prototype) error {
    if !open {
        return noConn
    }

    fields := getPrototypeCols()
    query := fmt.Sprintf("UPDATE prototypes SET %s=? %s=? WHERE %s=%s",
            fields[2], fields[3], fields[1], prototype.Name)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(prototype.DZone, prototype.DOffset)

    return err
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
    query := fmt.Sprintf("INSERT INTO replicas (%s, %s) VALUES(?, ?)", fields[1], fields[2])

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
    query := fmt.Sprintf("UPDATE userinfo SET %s=? WHERE %s=%s", fields[2], fields[1], replicaName)

    stmt, err := db.Prepare(query)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    return err
}

// AddZones adds a new sub-table of zones, to the specified prototype timezone.
func AddZones (prototypeName string, zones []Zone) error {
    if !open {
        return noConn
    }

/*
    prototypeID, err := getPrototypeID (prototypeName)
    prototype, err := getPrototype (prototypeID)

    // get name of most recent zone-table
    // compare table contents with zones slice
    // decide wht to do...

    err = createZones (newTableName, zones)
*/

    return nil
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
