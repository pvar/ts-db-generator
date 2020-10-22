package main

import (
        "fmt"
        "time"
)

type zoneTrans struct {
        when         int64 // transition time, in seconds since 1970 GMT
        index        uint8 // the index of the zone that goes into effect at that time
        isstd, isutc bool  // ignored - no idea what these mean
}

func main() {

        for _, location := range []string{"Europe/Berlin", "Europe/Athens"} {

                fmt.Println("\n\n" + location)

                loc, err := time.LoadLocation(location)
                if err != nil {
                        fmt.Println("Error:", err)
                        continue
                }

                now := time.Now().In(loc)
                zoneNameCurrent, timeOffset := now.Zone()
                zoneNameWinter, winterOffset := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc).Zone()
                zoneNameSummer, summerOffset := time.Date(now.Year(), 7, 1, 0, 0, 0, 0, loc).Zone()

                if winterOffset > summerOffset {
                        winterOffset, summerOffset = summerOffset, winterOffset
                        zoneNameWinter, zoneNameSummer = zoneNameSummer, zoneNameWinter
                }

                fmt.Println("current era:", zoneNameCurrent, timeOffset)
                fmt.Println("winter era: ", zoneNameWinter, winterOffset)
                fmt.Println("summer era: ", zoneNameSummer, summerOffset)

                var DS int;
                if winterOffset != summerOffset { // the location has Daylight Savings
                        if timeOffset != winterOffset {
                                DS = 1
                        } else {
                                DS = 0
                        }
                } else {
                        // the location does not have Daylight Savings
                        DS = -1
                }

                fmt.Println("offset: ", timeOffset)
                fmt.Println("DST: ", DS)
        }
}
