package tzdbio

import (
        "fmt"
        "testing"
)

func TestAddOriginal(t *testing.T) {

        originals := []string{
                "Shanghai",
                "Europe/Athens",
                "Europe/Berlin",
                "Portugal"}

        Open ("./testdb.sqlite")
        for _, original := range originals {
                id, err := AddOriginal (original)
                if err != nil {
                        fmt.Printf("Failed to add original %q\n", original)
                        t.Errorf("%s", err)
                } else {
                        fmt.Printf("Original %q added with ID %v\n", original, id)
                }
        }
}
