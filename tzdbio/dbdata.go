package tzdbio

import "fmt"

// Prototype defines a unique timezone
// By unique, we mean that said timezone is not a link to
// another one. It's kind of a prototype. All the zones
// withing a specific timezone are defined in a separate
// table.
type Prototype struct {
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
    prototypeTable string = "prototype"
    replicaTable string = "replica"
)

// column names for table of prototypes
func getPrototypeCols() []string {
    return []string{
        "id",
        "prototype_name",
        "default_zone",
        "default_offset",
        "ztable_name",
        "ztable_ver",
        "tzdada_ver"}
}

// column names for table of replicas
func getReplicaCols() []string {
    return []string{
        "id",
        "replica_name",
        "prototype_id"}
}

// column names for each tables of zones
func getZoneCols() []string {
    return []string{
        "id",
        "zone_abbrev",
        "zone_start",
        "zone_end",
        "zone_offset",
        "is_dst"}
}

// column names for table of prototypes
func getPrototypeSchema() string {
    fields := getPrototypeCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL, %q TEXT, %q INTEGER, %q TEXT, %q INTEGER DEFAULT 0, %q TEXT NOT NULL, PRIMARY KEY(%q AUTOINCREMENT));",
                        prototypeTable, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[6], fields[0])

    return schema
}

// column names for table of replicas
func getReplicaSchema() string {
    fields := getReplicaCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL, %q INTEGER NOT NULL, PRIMARY KEY(%q AUTOINCREMENT));",
                        replicaTable, fields[0], fields[1], fields[2], fields[0])

    return schema
}

// column names for each tables of zones
func getZoneSchema(name string) string {
    fields := getZoneCols()

    schema := fmt.Sprintf("CREATE TABLE %q (%q INTEGER UNIQUE, %q TEXT NOT NULL, %q INTEGER, %q INTEGER, %q INTEGER NOT NULL, %q INTEGER, PRIMARY KEY(%q AUTOINCREMENT));",
                        name, fields[0], fields[1], fields[2], fields[3], fields[4], fields[5], fields[0])

    return schema
}
