package tzdbio

import (
        "fmt"
        "testing"
)

func TestAddOriginal(t *testing.T) {

    originals := []string{
                "Shanghai",
                "Europe/Athens",
                "Europe/Berlin"}

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

func TestAddReplica(t *testing.T) {
        replicas := []struct{replica, original string} {
                {"Beijing", "Shanghai"},
                {"Shenzhen", "Shanghai"},
                {"Tianjin", "Shanghai"},
                {"Larisa", "Europe/Athens"},
                {"Farsala", "Europe/Athens"},
                {"Xalkida", "Europe/Athens"},
                {"Frankfurt", "Europe/Berlin"},
                {"Nurnberg", "Europe/Berlin"},
                {"Hanover", "Europe/Berlin"}}

        Open ("./testdb.sqlite")
        for _, replica := range replicas {
                err := AddReplicas([]string{replica.replica}, replica.original)
                if err != nil {
                        fmt.Printf("Failed to add replica %q with original %q\n", replica.replica, replica.original)
                        t.Errorf("%s", err)
                } else {
                        fmt.Printf("Added replica %q with original %q\n", replica.replica, replica.original)
                }
        }
}
