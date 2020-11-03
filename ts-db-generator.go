package main

import (
    "os"
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

    var filename string
    if len(os.Args) > 1 {
        filename = os.Args[1]
    } else {
        filename = dbfile
    }

    tzdb.Open(filename)
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

    storedCount, err := tzdb.GetOriginalCount()
    if err != nil || storedCount == 0 {
        // if no original timezones present,
        // assign a value that will in effect
        // disable the following check...
        storedCount = 123456789
    }

    // check if ammount of new originals supersedes 5% of stored originals
    if (float64(len(originals) - storedCount) / float64(storedCount)) > 0.05 {
        log.Printf("\nUpdated set of originals contains too many new entries!")
        return fmt.Errorf("")
    }

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

    storedCount, err := tzdb.GetReplicaCount()
    if err != nil || storedCount == 0 {
        // if no original timezones present,
        // assign a value that will in effect
        // disable the following check...
        storedCount = 123456789
    }

    // check if ammount of new replicas supersedes 5% of stored replicas
    if (float64(len(replicas) - storedCount) / float64(storedCount)) > 0.05 {
        log.Printf("\nUpdated set of originals contains too many new entries!")
        return fmt.Errorf("")
    }

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

    // loop through original timezones...
    i, j := 0, len(originals)
    for org, _ := range originals {
        i++
        // restore cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding full data of original timezone [%3d/%3d]", i, j)

        // get data related to selected timezone
        data, err := tzdata.GetData(org)
        if err != nil {
            log.Printf("\nfailed to get data for timezone %q: %s", org,  err)
            return err
        }

        originals[org].TZDVer = ver

        // These are the defualt values for Zone name (abbreviation)
        // and offset. They will be ignored if there are any zones
        // defined...
        zoneName, offset, _, _ := data.Lookup(nowTime)
        originals[org].DZone = zoneName
        originals[org].DOffset = int64(offset)

        // get metadata for already stored table of zones
        curTableVer, storedZones, storedTZdataVer, err := tzdb.GetZoneTableMeta (int(originals[org].ID))
        if err != nil {
            // if no stored zones are present,
            // assign a value that will in effect
            // disable the following check...
            storedZones = 123456789
        }

        saveZones := false
        zoneCount := len(data.Trans)

        // If frershly parsed data are of an older version than the stored data,
        // abort the update proceedure immediately.
        if (ver < storedTZdataVer) {
            return fmt.Errorf("Parsed TZdata are of an older version!")
        }

        // If freshly parsed and stored data are of the same version
        // AND the ammount of new zones equals the ammount stored ones,
        // there is nothing new to add.
        if (ver == storedTZdataVer) && (zoneCount == storedZones) {
            // Nothing new to add!
            // Proceed to next original timezone.
            continue
        }

        // Check if ammount of new zones supersedes 5% of ammount of stored zones.
        if ((float64(zoneCount - storedZones) / float64(storedZones)) > 0.05) {
            // Updated set of zones contains too many new entries!
            // Proceed to next original timezone.
            continue
        }

        // if any zones are defined, populate slice of zones
        zones := make([]tzdb.Zone, 0, zoneCount)
        if zoneCount != 0 {
            saveZones = true
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

            originals[org].TabVer = int64(curTableVer + 1)
            // name of zones-table will be auto-generated
        }

        // Save updated state of original timezone.
        if err := tzdb.UpdateOriginal(originals[org]); err != nil {
            log.Printf("\nattempt to update original %q failed with: %s", org,  err)
            return err
        }

        // If new zones are to be saved... do it!
        if saveZones {
            if err := tzdb.AddZones (org, zones); err != nil {
                log.Printf("\nattempt to add zones for original %q failed with: %s", org,  err)
                return err
            }
        }
    }
    return nil
}
