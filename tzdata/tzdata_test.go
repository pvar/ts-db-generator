package tzdata

import (
        "fmt"
        "testing"
)

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
                fmt.Printf("\nAttempt to get data for: %-24q Expecting error... ", badLocation)
                _, err := GetData(badLocation)
                if err == nil {
                        fmt.Printf("no error!!")
                        t.Errorf("Did not get error for %q!", badLocation)
                } else {
                        fmt.Printf("error occured")
                }
        }
        fmt.Println("")
}

func TestWierdLocations(t *testing.T) {
        badLocations := []string{
                "Asia/Shanghai",  // no DST
                "Etc/GMT",        // no transitions
                "Etc/GMT-14",     // large offset
                "right/Portugal"} // has leap seconds

        for _, badLocation := range badLocations {
                fmt.Printf("\nAttempt to get data for: %-24q Should get no errors... ", badLocation)
                _, err := GetData(badLocation)
                if err != nil {
                        fmt.Printf("got error!!")
                        t.Errorf("Got error for %q: %s", badLocation, err)
                } else {
                        fmt.Printf("no error")
                }
        }
        fmt.Println("")
}

func TestMalformedTZData(t *testing.T) {
        // The goal here is just that malformed tzdata results in an error, not a panic.
        issue29437 := "TZif\x00000000000000000\x00\x00\x00\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0000"
        _, err := parseRawTZdata("abc", []byte(issue29437))
        if err == nil {
                t.Error("expected error, got none")
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
