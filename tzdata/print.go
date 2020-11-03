package tzdata

import (
        "fmt"
        "time"
)

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
                var name string
                if trans.Index != 255 {
                    name = tzd.Eras[trans.Index].Name
                } else {
                    name = trans.AltName
                }
                fmt.Printf("        [%v] era: (%v) %-6s unix time: %-12v {isstd: %v, isutc: %v}\n",
                        i,
                        trans.Index,
                        name,
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
