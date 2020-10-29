package tzdbio

import (
        "fmt"
        "testing"
)

func init() {
        fmt.Printf("\nStarting tests...\n")
        Open ("./testdb.sqlite")
        query := fmt.Sprintf("DELETE FROM %s", originalTable)
        stmt, _ := db.Prepare(query)
        stmt.Exec()
        stmt.Close()
        query = fmt.Sprintf("DELETE FROM %s", replicaTable)
        stmt, _ = db.Prepare(query)
        stmt.Exec()
        stmt.Close()
        query = fmt.Sprintf("UPDATE sqlite_sequence SET seq=0 WHERE name=%q", originalTable)
        stmt, _ = db.Prepare(query)
        stmt.Exec()
        stmt.Close()
}

func TestAddOriginal(t *testing.T) {
    originals := []string{
                "Lamia",
                "Patra",
                "Irakleio"}

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

func TestUpdateOriginal(t *testing.T) {
    originals := []Original{
                {ID: 0, Name: "Lamia", DZone: "kalokairi", DOffset: 213, TabName: "lamia_test", TabVer: 0, TZDVer: "2020a"},
                {ID: 0, Name: "Patra", DZone: "kalokairi", DOffset: 123, TabName: "patra_test", TabVer: 0, TZDVer: "2020a"},
                {ID: 0, Name: "Irakleio", DZone: "kalokairi", DOffset: 312, TabName: "irkleio_test", TabVer: 0, TZDVer: "2020a"}}

        for _, original := range originals {
                err := UpdateOriginal (&original)
                if err != nil {
                        fmt.Printf("Failed to update data for original %q\n", original)
                        t.Errorf("%s", err)
                } else {
                        fmt.Printf("Updated data for original %q\n", original.Name)
                }
        }
}

func TestAddZones(t *testing.T) {
        var original = "Lamia"
        var zones = []Zone{
        {ID: 1, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 2, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 3, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 4, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 5, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 6, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 7, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 8, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 9, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 10, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 11, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 12, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 13, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 14, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false},
        {ID: 15, Name: "kalokairi", Start: 1234, End: 5678, Offset: 1589, IsDST: true},
        {ID: 16, Name: "xeimonas", Start: 5678, End: 1234, Offset: 1581, IsDST: false}}

        err := AddZones(original, zones)
        if err != nil {
                fmt.Printf("Failed to add zones to original %q\n", original)
                t.Errorf("%s", err)
        } else {
                fmt.Printf("Added %d zones to original %q\n", len(zones), original)
        }
}

func TestAddReplica(t *testing.T) {
        replicas := []struct{replica, original string} {
                {"Lamia", "Lamia"},
                {"Orxomenos", "Lamia"},
                {"Lianokladi", "Lamia"},
                {"Patra", "Patra"},
                {"Leontio", "Patra"},
                {"Antirio", "Patra"},
                {"Iraklio", "Irakleio"},
                {"Knosos", "Irakleio"},
                {"Finikia", "Irakleio"}}

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

func BenchmarkGetOriginalByID(b *testing.B) {
        var testID = 3
        for i := 0; i < b.N; i++ {
                _, err := getOriginalByID (testID)
                if err != nil {
                        b.Errorf("%s", err)
                }
        }
}

func BenchmarkGetOriginalByNane(b *testing.B) {
        for i := 0; i < b.N; i++ {
                _, err := getOriginalByName ("Lamia")
                if err != nil {
                        b.Errorf("%s", err)
                }
        }
}

func BenchmarkGetReplicaOriginal(b *testing.B) {
        for i := 0; i < b.N; i++ {
                _, err := getReplicaOriginal ("Orxomenos")
                if err != nil {
                        b.Errorf("%s", err)
                }
        }
}

func BenchmarkGetZones(b *testing.B) {
        var original = "Lamia"
        for i := 0; i < b.N; i++ {
                _, err := GetZones (original)
                if err != nil {
                        b.Errorf("%s", err)
                }
        }
}
