// tzdata parses zones and transitions from timezone files.
// tzdata is actually the part of time package that handles
// locations. The code is stripped down to the absolute
// minimum, in order to only run on Linux and always use
// the timezone files installed on the system. All available
// data are exposed and some methods for printing data
// in a meaningful(*) way have been added.
//
// (*) Yes, tzdata was built with a specific application in mind
// and it is doubtful it will be of any use to others.

package tzdata

import (
        "fmt"
        "time"
)

// TZdata collects time offsets and offset-transitions for a geographical area.
// Typically, the TZdata struct represents the collection of time offsets
// in use in a geographical area. For many Locations the time offset varies
// depending on whether daylight savings time is in use.
type TZdata struct {
        Name  string
        Eras  []Era
        Trans []EraTrans

        // The tzdata information can be followed by a string that describes
        // how to handle DST transitions not recorded in zoneTrans.
        // The format is the TZ environment variable without a colon;
        // https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html.
        Extend string
}

// A zone represents a single time zone (CET, CEST, etc).
type Era struct {
        Name   string // abbreviated name of zone
        Offset int    // seconds east of UTC
        IsDST  bool   // is this zone Daylight Savings Time?
}

// A zoneTrans represents a single time zone transition.
type EraTrans struct {
        When         int64 // transition time, in seconds since 1970 GMT
        Index        uint8 // index of the zone that goes into effect at that time
        Isstd, Isutc bool  // seems to be ignored
        // supposed to indicate whether transition time (When)
        // expressed in UTC or local
}

func (tzd *TZdata) PrintRaw() {
        fmt.Printf("\nRaw data for %q.\n", tzd.Name)

        fmt.Printf("    era names:\n")
        for i, era := range tzd.Eras {
                fmt.Printf("        [%v] name: %-5s offset: %-6v DST: %v\n",
                        i,
                        era.Name,
                        era.Offset,
                        era.IsDST)
        }
        fmt.Printf("    transitions:\n")
        for i, trans := range tzd.Trans {
                fmt.Printf("        [%v] era: (%v) %-6s unix time: %-12v {isstd: %v, isutc: %v}\n",
                        i,
                        trans.Index,
                        tzd.Eras[trans.Index].Name,
                        trans.When,
                        trans.Isstd,
                        trans.Isutc)
        }
        fmt.Printf("    TZ variable: %s\n", tzd.Extend)
}

func (tzd *TZdata) PrintProcessed() {
        fmt.Printf("\nProcessed data for %s:\n", tzd.Name)

        now := time.Now()
        nowEpoch := now.Unix()

        name, offset, start, end := tzd.Lookup(nowEpoch)
        fmt.Printf("    Current Era: %q\n", name)
        fmt.Printf("        Offset : %v\n", offset)

        var t time.Time
        if start == bigbang  || start == gnabgib {
                fmt.Printf("        start  : bigbang\n")
        } else {
                t = time.Unix(start, 0)
                fmt.Printf("        start  : %s -- %v\n", t.Format(time.ANSIC), start)
        }

        if end == gnabgib {
                fmt.Printf("        stop   : gnabgib\n")
        } else {
                t = time.Unix(end, 0)
                fmt.Printf("        stop   : %s -- %v\n", t.Format(time.ANSIC), end)
        }

        if end < gnabgib {
                name, offset, start, end = tzd.Lookup(end+1)
                fmt.Printf("    Comming Era: %q\n", name)
                fmt.Printf("        Offset : %v\n", offset)

                var t time.Time
                if start == bigbang  || start == gnabgib {
                        fmt.Printf("        start  : bigbang\n")
                } else {
                        t = time.Unix(start, 0)
                        fmt.Printf("        start  : %s -- %v\n", t.Format(time.ANSIC), start)
                }

                if end == gnabgib {
                        fmt.Printf("        stop   : gnabgib\n")
                } else {
                        t = time.Unix(end, 0)
                        fmt.Printf("        stop   : %s -- %v\n", t.Format(time.ANSIC), end)
                }
        }
}
