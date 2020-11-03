package tzdata

import (
        "fmt"
        "testing"
)

func TestGetData(t *testing.T) {
    location := "Europe/Athens"

    data, err := GetData(location)
    if err != nil {
        t.Errorf("Error getting data for %q: %s", location, err)
    }

    fmt.Printf("\nTimezone: %q", data.Name)

    for i, era := range data.Eras {
        fmt.Printf("\n\t[%3d] Era: %-5q, DST: %v, Offset: %d", i, era.Name, era.IsDST, era.Offset)
    }
    fmt.Printf("\n")

    for i, trans := range data.Trans {
        if trans.Index != 255 {
            fmt.Printf("\n\t[%3d] Name: %6q\tStart: %d\tOffset: %d", i, data.Eras[trans.Index].Name, trans.When, data.Eras[trans.Index].Offset)
        } else {
            fmt.Printf("\n\t[%3d] Name: !%5q\tStart: %d\tOffset: %d", i, trans.AltName, trans.When, trans.AltOffset)
        }
    }

    fmt.Printf("\n")

}

func TestGetList(t *testing.T) {
    version, list, err := GetList()

    if err != nil {
        t.Errorf("\nFailed: %s\n", err)
    } else {
        fmt.Printf("\n\tVersion: %q\n", version)
        fmt.Printf("\tTimezones:\n")
        for key, value := range list {
            fmt.Printf("\t\t%q --> %q\n", key, value)
        }
    }
}

func BenchmarkGetList(b *testing.B) {
        for i := 0; i < b.N; i++ {
            GetList()
        }
}

func TestBadLocations(t *testing.T) {
        badLocations := []string{
                "Atlantis/nonexistent",
                "/Asia/Shanghai",
                "Asia\\Shanghai",
                "..Europe/Athens",
                "Athens",
                "Europe",
                "..",
                "/",
                ""}

        for _, badLocation := range badLocations {
                _, err := GetData(badLocation)
                if err == nil {
                        fmt.Printf("\nAttempt to get data for %-24q should produce an error, but it did not!", badLocation)
                        t.Errorf("Did not get error for %q!", badLocation)
                }
        }
}

func TestWierdLocations(t *testing.T) {
        badLocations := []string{
                "Asia/Shanghai",  // no DST
                "Etc/GMT",        // no transitions
                "Etc/GMT-14",     // large offset
                "right/Portugal"} // has leap seconds

        for _, badLocation := range badLocations {
                _, err := GetData(badLocation)
                if err != nil {
                        fmt.Printf("\nAttempt to get data for %-24q should get no errors, but it did!", badLocation)
                        t.Errorf("Got error for %q: %s", badLocation, err)
                }
        }
}

func TestMalformedTZData(t *testing.T) {
        // The goal here is just that malformed tzdata results in an error, not a panic.
        issue29437 := "TZif\x00000000000000000\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0000"
        _, err := parseRawTZdata("abc", []byte(issue29437))
        if err == nil {
                t.Error("expected error, got none")
        }
}

func TestTzset(t *testing.T) {
        for _, test := range []struct {
                inStr string
                inEnd int64
                inSec int64
                name  string
                off   int
                start int64
                end   int64
                ok    bool
        }{
                {"", 0, 0, "", 0, 0, 0, false},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2159200800, "PDT", -7 * 60 * 60, 2152173600, 2172733200, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2152173599, "PST", -8 * 60 * 60, 2145916800, 2152173600, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2152173600, "PDT", -7 * 60 * 60, 2152173600, 2172733200, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2152173601, "PDT", -7 * 60 * 60, 2152173600, 2172733200, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2172733199, "PDT", -7 * 60 * 60, 2152173600, 2172733200, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2172733200, "PST", -8 * 60 * 60, 2172733200, 2177452800, true},
                {"PST8PDT,M3.2.0,M11.1.0", 0, 2172733201, "PST", -8 * 60 * 60, 2172733200, 2177452800, true},
        } {
                name, off, start, end, ok := tzset(test.inStr, test.inEnd, test.inSec)
                if name != test.name || off != test.off || start != test.start || end != test.end || ok != test.ok {
                        t.Errorf("tzset(%q, %d, %d) = %q, %d, %d, %d, %t, want %q, %d, %d, %d, %t", test.inStr, test.inEnd, test.inSec, name, off, start, end, ok, test.name, test.off, test.start, test.end, test.ok)
                }
        }
}

func TestTzsetName(t *testing.T) {
        for _, test := range []struct {
                in   string
                name string
                out  string
                ok   bool
        }{
                {"", "", "", false},
                {"X", "", "", false},
                {"PST", "PST", "", true},
                {"PST8PDT", "PST", "8PDT", true},
                {"PST-08", "PST", "-08", true},
                {"<A+B>+08", "A+B", "+08", true},
        } {
                name, out, ok := tzsetName(test.in)
                if name != test.name || out != test.out || ok != test.ok {
                        t.Errorf("tzsetName(%q) = %q, %q, %t, want %q, %q, %t", test.in, name, out, ok, test.name, test.out, test.ok)
                }
        }
}

func TestTzsetOffset(t *testing.T) {
        for _, test := range []struct {
                in  string
                off int
                out string
                ok  bool
        }{
                {"", 0, "", false},
                {"X", 0, "", false},
                {"+", 0, "", false},
                {"+08", 8 * 60 * 60, "", true},
                {"-01:02:03", -1*60*60 - 2*60 - 3, "", true},
                {"01", 1 * 60 * 60, "", true},
                {"100", 0, "", false},
                {"8PDT", 8 * 60 * 60, "PDT", true},
        } {
                off, out, ok := tzsetOffset(test.in)
                if off != test.off || out != test.out || ok != test.ok {
                        t.Errorf("tzsetName(%q) = %d, %q, %t, want %d, %q, %t", test.in, off, out, ok, test.off, test.out, test.ok)
                }
        }
}

func BenchmarkLargeTZfile(b *testing.B) {
        for i := 0; i < b.N; i++ {
                GetData("Europe/Belfast")
        }
}

func BenchmarkSmallTZfile(b *testing.B) {
        for i := 0; i < b.N; i++ {
                GetData("Zulu")
        }
}
