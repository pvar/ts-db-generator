package main

import (
    "fmt"
    "log"
    "time"
    "ts-db-generator/tzdata"
    "ts-db-generator/tzdb"
)

const dbfile = "./tsdb.sqlite"

func main () {
    version, timezones, err := tzdata.GetList()
    if err != nil {
        log.Fatalf("\nError loading timezone metadata (tzdata.zi): %s", err)
    }

    originals := make(map[string]*tzdb.Original)
    replicas := make(map[string][]string)
    for replica, original := range timezones {
        if originals[original] == nil {
            originals[original] = &tzdb.Original{Name: original}
        }
        replicas[original] = append(replicas[original], replica)
    }

    tzdb.Open(dbfile)
    defer tzdb.Close()

    if err := storeOriginals(originals); err != nil {
        log.Fatalf("\nFailed while storing originals")
    }

    if err := storeReplicas(replicas); err != nil {
        log.Fatalf("\nFailed while storing replicas")
    }

    if err := updateOriginals (version, originals); err != nil {
        log.Fatalf("\nFailed while updating originals")
    }

    fmt.Printf("\nAll done. Have a nice day :)\n")
}

// storeOriginals add new entries in the table of original timezones
// THe ID of each entry is saved in the struct representing each
// timezone, since it will be needed later-on, while storing the
// replicas (links to originals).
func storeOriginals (originals map[string]*tzdb.Original) error {
    // save cursor position
    fmt.Print("\033[s")

    i, j := 0, len(originals)
    for org, _ := range originals {
        i++
        // restore cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding original timezone [%3d/%3d]", i, j)

        id, err := tzdb.AddOriginal (org)
        if err != nil {
            log.Printf("\nattempt to add %q failed with: %s", org,  err)
            return err
        }
        originals[org].ID = id
    }
    fmt.Print("\n")
    return nil
}

// storeReplicas stores groups of replica-timezones.
// That is, timezones that are linked to another timezone
// and refer to the same set of data.
func storeReplicas(replicas map[string][]string) error {
    // save cursor position
    fmt.Print("\033[s")

    i, j := 0, len(replicas)
    for org, rlist := range replicas {
        i++
        // restore cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding group of replicas [%3d/%3d]", i, j)

        err := tzdb.AddReplicas (rlist, org)
        if err != nil {
            log.Printf("\nattempt to add %q failed with: %s", org,  err)
            return err
        }
    }
    fmt.Print("\n")
    return nil
}

// updateOriginals stores all related to each original timezone.
// That is, all the available zones, the default zone and offset
// and the version of the tzdata set used.
func updateOriginals (ver string, originals map[string]*tzdb.Original) error {
    // save cursor position
    fmt.Print("\033[s")

    nowTime := time.Now().Unix()

    i, j := 0, len(originals)
    for org, _ := range originals {
        i++
        // restore cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding full data of original timezone [%3d/%3d]", i, j)

        originals[org].TZDVer = ver

        // get data related to selected timezone
        data, err := tzdata.GetData(org)
        if err != nil {
            log.Printf("\nattempt to get data for original %q failed with: %s", org,  err)
            return err
        }

        // these values will be ignored if there are any zones defined
        zoneName, offset, _, _ := data.Lookup(nowTime)
        originals[org].DZone = zoneName
        originals[org].DOffset = int64(offset)

        // check if any zones are defined
        saveZones := false
        zoneCount := len(data.Trans)
        zones := make([]tzdb.Zone, zoneCount)
        if zoneCount != 0 {
            zone := tzdb.Zone{}
            for z := 0; z < zoneCount; z++ {
                zone.Name = data.Eras[data.Trans[z].Index].Name
                zone.Offset = int64(data.Eras[data.Trans[z].Index].Offset)
                zone.IsDST = data.Eras[data.Trans[z].Index].IsDST
                zone.Start = data.Trans[z].When
                if z+1 < zoneCount {
                    zone.End = data.Trans[z+1].When - 1
                } else {
                    zone.End = -1 // end of time!
                }
                zones = append(zones, zone)
            }

            // get current version of zones table
            //...

            // update version of zones-table
            // name of zones-table will be auto-generated
            originals[org].TabVer += 1

            saveZones = true
        }

        // save current state of original
        if err := tzdb.UpdateOriginal(originals[org]); err != nil {
            log.Printf("\nattempt to update original %q failed with: %s", org,  err)
            return err
        }

        // if have to, save table with zones
        if saveZones {
            if err := tzdb.AddZones (org, zones); err != nil {
                log.Printf("\nattempt to add zones for original %q failed with: %s", org,  err)
                return err
            }
        }
    }
    return nil
}
