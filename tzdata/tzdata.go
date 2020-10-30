// tzdata parses zones and transitions from timezone files.
// tzdata is actually the part of time package that handles
// locations. The code is stripped down to the absolute
// minimum, in order to only run on Linux and always use
// the timezone files installed on the system. All available
// data are exposed and some methods for printing data
// in a meaningful(*) way have been added.
//
// (*) Yes, tzdata was built with a specific application in mind
// and it is doubtful it will be of any use to others.

package tzdata

const source_path string = "/usr/share/zoneinfo/"

// TZdata collects time offsets and offset-transitions for a geographical area.
// Typically, the TZdata struct represents the collection of time offsets
// in use in a geographical area. For many Locations the time offset varies
// depending on whether daylight savings time is in use.
type TZdata struct {
        Name  string
        Eras  []Era
        Trans []EraTrans

        // The tzdata information can be followed by a string that describes
        // how to handle DST transitions not recorded in zoneTrans.
        // The format is the TZ environment variable without a colon;
        // https://pubs.opengroup.org/onlinepubs/9699919799/basedefs/V1_chap08.html.
        Extend string
}

// A zone represents a single time zone (CET, CEST, etc).
type Era struct {
        Name   string // abbreviated name of zone
        Offset int    // seconds east of UTC
        IsDST  bool   // is this zone Daylight Savings Time?
}

// A zoneTrans represents a single time zone transition.
type EraTrans struct {
        When         int64 // transition time, in seconds since 1970 GMT
        Index        uint8 // index of the zone that goes into effect at that time
        Isstd, Isutc bool  // seems to be ignored
        // supposed to indicate whether transition time (When)
        // expressed in UTC or local
}
