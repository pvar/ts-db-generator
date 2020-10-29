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
                        fmt.Printf("Added original %q with ID %v\n", original, id)
                }
        }
}

func TestAddFullOriginal(t *testing.T) {
    originals := []Original{
                {ID: 0, Name: "Lamia", DZone: "kalokairi", DOffset: 213, TabName: "lamia_test", TabVer: 0, TZDVer: "2020a"},
                {ID: 0, Name: "Patra", DZone: "kalokairi", DOffset: 123, TabName: "patra_test", TabVer: 0, TZDVer: "2020a"},
                {ID: 0, Name: "Irakleio", DZone: "kalokairi", DOffset: 312, TabName: "irkleio_test", TabVer: 0, TZDVer: "2020a"}}

        Open ("./testdb.sqlite")
        for _, original := range originals {
                err := AddFullOriginal (&original)
                if err != nil {
                        fmt.Printf("Failed to add full data for original %q\n", original)
                        t.Errorf("%s", err)
                } else {
                        fmt.Printf("Added full data for original %q\n", original.Name)
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
                {"Finikia", "Irakleio"}}

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
