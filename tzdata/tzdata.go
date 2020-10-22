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
    transitionsCount := len(tzd.Trans)

    // If no transitions defined,
    // consider the last entry in eras slice as the current era.
    if transitionsCount == 0 {
        fmt.Println("!! No transitions defined. Will pick last defined era...")
        fmt.Printf("    current era: %s\n", tzd.Eras[len(tzd.Eras)-1].Name)
        fmt.Printf("        offset : %v\n", tzd.Eras[len(tzd.Eras)-1].Offset)
        fmt.Printf("        end    : unknown\n")
        return
    }

    // Transitions are always sorted in ascending order...
    now := time.Now()
    nowEpoch := now.Unix()
    foundTransition := false
    nextTransIndex := 0
    for ; nextTransIndex < transitionsCount; nextTransIndex++ {
        if nowEpoch < tzd.Trans[nextTransIndex].When {
            foundTransition = true
            break
        }
    }

    // If no transitions ahead,
    // current era is defined by the last transition.
    if !foundTransition {
        fmt.Println("!! No transitions in foreseeable future. Last transition determines current era.")
        fmt.Printf("    current era: %s\n", tzd.Eras[tzd.Trans[transitionsCount-1].Index].Name)
        fmt.Printf("        offset : %v\n", tzd.Eras[tzd.Trans[transitionsCount-1].Index].Offset)
        fmt.Printf("        end    : gnab gib :-)\n")
        return
    }

    // The immediately previous transition specifies the name of current era.
    // If the next transition is the first in the slice, there is no certan way
    // to deduce the name of the current era.
    if nextTransIndex == 0 {
        fmt.Printf("!! No past transition found. Cannot deduce name of curent era. Aborting...")
        return
    }

    fmt.Printf("    current era: %s\n", tzd.Eras[tzd.Trans[nextTransIndex-1].Index].Name)
    fmt.Printf("        offset : %v\n", tzd.Eras[tzd.Trans[nextTransIndex-1].Index].Offset)
    end := time.Unix(tzd.Trans[nextTransIndex].When-1, 0)
    fmt.Printf("        end    : (%v) %s\n", tzd.Trans[nextTransIndex].When-1, end.Format("Mon, 02 01 2006 15:04:05"))

    fmt.Printf("    comming era: %s\n", tzd.Eras[tzd.Trans[nextTransIndex].Index].Name)
    fmt.Printf("        offset : %v\n", tzd.Eras[tzd.Trans[nextTransIndex].Index].Offset)
    start := time.Unix(tzd.Trans[nextTransIndex].When, 0)
    fmt.Printf("        start  : (%v) %s\n", tzd.Trans[nextTransIndex].When, start.Format("Mon, 02 01 2006 15:04:05"))
}
