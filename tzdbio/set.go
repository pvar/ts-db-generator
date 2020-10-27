package tzdbio

import (
    "fmt"
//    "errors"
)

// AddPrototype adds *only* the name of a new entry in prototypes' table.
// The rest of the data remain uninitialized. This function is mainly used
// during initial setup, to populate table with available prototypes.
func AddPrototype (prototypeName string) (id int64, err error) {
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

// AddReplicas adds a new list of entries in the preplicas' table.
// Each group of replicas contains the name of the original as an
// extra entry. Mainly used during initial setup, to populate table
// with replicas.
func AddReplicas (replicas []string, prototypeName string) error {
    return nil
}

//
//
//
func UpdateReplicas (replicaName, prototypeName string) error {
    if !open {
        return noConn
    }

    prototypeID, err := getPrototypeByName (prototypeName)
    if err != nil {
        // specified prototype cannot be found, will attempt to add it...
        // the following prototype is hollow -- it lacks all useful data
        newPrototype := Prototype{Name: prototypeName}
        err := AddPrototype (newPrototype)
        if err != nil {
            // failed trying to append specified prototype
            return err
        }
    }

    err = updateReplica (replicaName, prototypeID)
    return err
}

//
//
// Used to add a new set of zones to an existing prototype.
func AddZonesTable (prototypeName string, zones []Zone) error {
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




func updatePrototype (prototype Prototype) error {
    if !open {
        return noConn
    }

    return nil
}

func createZones (updatedTableName string, zones []Zone) error {
    return nil
}

// update link of replica to prototype
func updateReplica (name string, prototypeID int) error {
    return nil
}
