package main

import (
    "fmt"
    "log"
    "time"
    "ts-db-generator/tzdata"
    "ts-db-generator/tzdbio"
)

const dbfile = "./tsdb.sqlite"

func main () {
    version, timezones, err := tzdata.GetList()
    if err != nil {
        log.Fatalf("\n%s", err)
    }

    originals := make(map[string]*tzdbio.Original)
    replicas := make(map[string][]string)

    for replica, original := range timezones {
        if originals[original] == nil {
            originals[original] = &tzdbio.Original{Name: original}
        }

        replicas[original] = append(replicas[original], replica)
    }

    tzdbio.Open(dbfile)
    defer tzdbio.Close()

    fmt.Printf("\nVersion: %q\n", version)

    // save cursor position
    fmt.Print("\033[s")

    var i, j int

    // add original timezones and get respective ID for each one
    i, j = 0, len(originals)
    for org, _ := range originals {
        i++
        // restore the cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding original timezone [%3d/%3d]", i, j)

        id, err := tzdbio.AddOriginal (org)
        if err != nil {
            log.Fatalf("\nattempt to add %q failed with: %s", org,  err)
        }
        originals[org].ID = id
    }

    // move and save cursor position
    fmt.Print("\n\033[s")

    // add replicas and respective original ID
    i, j = 0, len(replicas)
    for org, rlist := range replicas {
        i++
        // restore the cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding group of replicas [%3d/%3d]", i, j)

        err := tzdbio.AddReplicas (rlist, org)
        if err != nil {
            log.Fatalf("\nattempt to add %q failed with: %s", org,  err)
        }
    }

    // move and save cursor position
    fmt.Print("\n\033[s")

    // get current time(stamp)
    nowTime := time.Now().Unix()

    // add full data and zones for each original
    i, j = 0, len(originals)
    for org, _ := range originals {
        i++
        // restore the cursor position and clear line
        fmt.Print("\033[u\033[K")
        fmt.Printf("Adding full data of original timezone [%3d/%3d]", i, j)

        // start with an easy one: tzdata version
        originals[org].TZDVer = version

        // name of zones table will be auto-generated
        //originals[org].TabName = ...

        // get tz data for selected original timezone
        data, err := tzdata.GetData(org)
        if err != nil {
            log.Fatalf("\nattempt to get data for original %q failed with: %s", org,  err)
        }

        // get default zone name default offset
        // (these values will be ignored if there are any zones defined)
        zoneName, offset, _, _ := data.Lookup(nowTime)
        originals[org].DZone = zoneName
        originals[org].DOffset = int64(offset)

        // check number of zone-transitions
        zoneCount := len(data.Trans)
        if zoneCount == 0 {
            // no more data to add...
            // save current state of original
            err := tzdbio.UpdateOriginal(originals[org])
            if err != nil {
                log.Fatalf("\nattempt to update original %q failed with: %s", org,  err)
            }
            // and proceed to next original
            continue
        }

        // it seems there are some zones to be saved...
        zones := make([]tzdbio.Zone, zoneCount)
        zone := tzdbio.Zone{}

        for i = 0; i < zoneCount; i++ {
            zone.Name = data.Eras[data.Trans[i].Index].Name
            zone.Offset = int64(data.Eras[data.Trans[i].Index].Offset)
            zone.IsDST = data.Eras[data.Trans[i].Index].IsDST
            zone.Start = data.Trans[i].When
            if i+1 < zoneCount {
                zone.End = data.Trans[i+1].When - 1
            } else {
                zone.End = -1 // end of time!
            }
            zones = append(zones, zone)
        }

        // get current version of zones table
        //...

        // update version of zones table
        originals[org].TabVer += 1

        // save current state of original
        if err := tzdbio.UpdateOriginal(originals[org]); err != nil {
            log.Fatalf("\nattempt to update original %q failed with: %s", org,  err)
        }

        // save table with zones
        if err := tzdbio.AddZones (org, zones); err != nil {
            log.Fatalf("\nattempt to add zones for original %q failed with: %s", org,  err)
        }
    }
}
