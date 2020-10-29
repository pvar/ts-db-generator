package tzdbio

import "fmt"

// Original defines a unique timezone.
// One that is not a link to another one.
// All zones withing a specific timezone
// are defined in separate tables.
type Original struct {
    ID      int64
    Name    string
    DZone   string      // Default Zone (when no zones are defined)
    DOffset int64       // Default Offset (when no zones are defined)
    TabName string
    TabVer  int64
    TZDVer  string      // Version of TZ-data used to update sqlite database
}

// Replica defines a link to some timezone
type Replica struct {
    ID      int64
    Name    string
    ProtoID int64
}

// Zone defines a zone within a timezone
type Zone struct {
    ID      int64
    Name    string
    Start   int64
    End     int64
    Offset  int64
    IsDST   bool
}

const (
    originalTable string = "original"
    replicaTable string = "replica"
)

// column names for table of prototypes
func getOriginalCols() []string {
    return []string{
        "id",
        "name",
        "default_zone",
        "default_offset",
        "zones_tab_name",
        "zones_tab_ver",
        "tzdada_ver"}
}

// column names for table of replicas
func getReplicaCols() []string {
    return []string{
        "id",
        "name",
        "original_id"}
}

// column names for each tables of zones
func getZoneCols() []string {
    return []string{
        "id",
        "abbrev",
        "start",
        "end",
        "offset",
        "is_dst"}
}

// column names for table of prototypes
func getOriginalSchema() string {
    fields := getOriginalCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL UNIQUE, %q TEXT, %q INTEGER, %q TEXT, %q INTEGER DEFAULT 0, %q TEXT, PRIMARY KEY(%q AUTOINCREMENT));",
                        originalTable, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[0])

    return schema
}

// column names for table of replicas
func getReplicaSchema() string {
    fields := getReplicaCols()
    fgnfields := getOriginalCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL, %q INTEGER NOT NULL, PRIMARY KEY(%q AUTOINCREMENT), FOREIGN KEY(%q) REFERENCES %s(%q));",
                        replicaTable, fields[0], fields[1], fields[2], fields[0], fields[2], originalTable, fgnfields[0])

    return schema
}

// column names for each tables of zones
func getZoneSchema(name string) string {
    fields := getZoneCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL, %q INTEGER, %q INTEGER, %q INTEGER NOT NULL, %q INTEGER, PRIMARY KEY(%q AUTOINCREMENT));",
                        name, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[0])

    return schema
}
