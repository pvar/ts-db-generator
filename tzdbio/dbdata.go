package tzdbio

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

// column names for table of prototypes
func getPrototypeCols() []string {
    return []string{
        "id",
        "prototype_name",
        "default_zone",
        "default_offset",
        "ztable_name",
        "ztable_ver"}
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
