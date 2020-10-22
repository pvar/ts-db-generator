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
		fmt.Printf("\nAttempt to get data for: %-24q Expecting some error... ", badLocation)
		_, err := GetData(badLocation)
		if err == nil {
			t.Errorf("Did not get error for %q!", badLocation)
		} else {
			fmt.Printf("got it!")
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
		fmt.Printf("\nAttempt to get data for: %-24q Should not get any errors... ", badLocation)
		_, err := GetData(badLocation)
		if err != nil {
			t.Errorf("Got error for %q: %s", badLocation, err)
		} else {
			fmt.Printf("all's good!")
		}
	}
	fmt.Println("")
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
