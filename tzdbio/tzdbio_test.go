package tzdbio

import (
        "fmt"
        "testing"
)

func TestAddOriginal(t *testing.T) {

    originals := []string{
                "Lamia",
                "Patra",
                "Irakleio"}

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
                {"Lamia", "Lamia"},
                {"Orxomenos", "Lamia"},
                {"Lianokladi", "Lamia"},
                {"Patra", "Patra"},
                {"Leontio", "Patra"},
                {"Antririo", "Patra"},
                {"Iraklio", "Irakleio"},
                {"Knosos", "Irakleio"},
                {"Finikia", "Irakleio"},
                {"Lamia", "Irakleio"}}

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
