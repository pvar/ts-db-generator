package tzdata

import (
	"errors"
	"runtime"
	"time"
)

var badData = errors.New("tzdata: malformed timezone file")

// Simple I/O interface to binary blob of data.
type dataIO struct {
	p     []byte
	error bool
}

func (d *dataIO) read(n int) []byte {
	if len(d.p) < n {
		d.p = nil
		d.error = true
		return nil
	}
	p := d.p[0:n]
	d.p = d.p[n:]
	return p
}

func (d *dataIO) big4() (n uint32, ok bool) {
	p := d.read(4)
	if len(p) < 4 {
		d.error = true
		return 0, false
	}
	return uint32(p[3]) | uint32(p[2])<<8 | uint32(p[1])<<16 | uint32(p[0])<<24, true
}

func (d *dataIO) big8() (n uint64, ok bool) {
	n1, ok1 := d.big4()
	n2, ok2 := d.big4()
	if !ok1 || !ok2 {
		d.error = true
		return 0, false
	}
	return (uint64(n1) << 32) | uint64(n2), true
}

func (d *dataIO) byte() (n byte, ok bool) {
	p := d.read(1)
	if len(p) < 1 {
		d.error = true
		return 0, false
	}
	return p[0], true
}

func (d *dataIO) rest() []byte {
	r := d.p
	d.p = nil
	return r
}

// Make a string by stopping at the first NUL
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// parseRawTZdata returns a TZdata with the given name initialized
// from the IANA Time Zone database-formatted data. The data should
// be in the format of a standard IANA time zone file.
func parseRawTZdata(name string, data []byte) (*TZdata, error) {
	d := dataIO{data, false}

	// 4-byte magic "TZif"
	if magic := d.read(4); string(magic) != "TZif" {
		return nil, badData
	}

	// 1-byte version, then 15 bytes of padding
	var version int
	var p []byte
	if p = d.read(16); len(p) != 16 {
		return nil, badData
	} else {
		switch p[0] {
		case 0:
			version = 1
		case '2':
			version = 2
		case '3':
			version = 3
		default:
			return nil, badData
		}
	}

	// six big-endian 32-bit integers:
	//  number of UTC/local indicators
	//  number of standard/wall indicators
	//  number of leap seconds
	//  number of transition times
	//  number of local time zones
	//  number of characters of time zone abbrev strings
	const (
		NUTCLocal = iota // isutcnt
		NStdWall         // isstdcnt
		NLeap            // leapcnt
		NTime            // timecnt
		NZone            // typecnt
		NChar            // charcnt
	)
	var cnt [6]int
	for i := 0; i < 6; i++ {
		cntval, ok := d.big4()
		if !ok {
			return nil, badData
		}
		if uint32(int(cntval)) != cntval {
			return nil, badData
		}
		cnt[i] = int(cntval)
	}

	// If we have version 2 or 3, then the data is first written out
	// in a 32-bit format, then written out again in a 64-bit format.
	// Skip the 32-bit format and read the 64-bit one, as it can
	// describe a broader range of dates.

	is64 := false
	if version > 1 {
		// Skip the 32-bit data (version 1 data block)
		// See more on [page 8] of RFC8536.
		skip := cnt[NTime]*4 +
			cnt[NTime] +
			cnt[NZone]*6 +
			cnt[NChar] +
			cnt[NLeap]*8 +
			cnt[NStdWall] +
			cnt[NUTCLocal]
		// Skip the first part (magic number, version and padding) of version 2 header.
		skip += 4 + 16
		d.read(skip)

		is64 = true

		// Read the counts again, they can differ.
		for i := 0; i < 6; i++ {
			cntval, ok := d.big4()
			if !ok {
				return nil, badData
			}
			if uint32(int(cntval)) != cntval {
				return nil, badData
			}
			cnt[i] = int(cntval)
		}
	}

	size := 4
	if is64 {
		size = 8
	}

	// Read data block (see more on [page 8] of RFC8536).

	// Transition times.
	txtimes := dataIO{d.read(cnt[NTime] * size), false}

	// Time zone indices for transition times.
	txzones := d.read(cnt[NTime])

	// Zone info structures.
	zonedata := dataIO{d.read(cnt[NZone] * 6), false}

	// Time zone abbreviations.
	abbrev := d.read(cnt[NChar])

	// Leap-second time pairs
	d.read(cnt[NLeap] * (size + 4))

	// Whether tx times associated with local time types
	// are specified as standard time or wall time.
	isstd := d.read(cnt[NStdWall])

	// Whether tx times associated with local time types
	// are specified as UTC or local time.
	isutc := d.read(cnt[NUTCLocal])

	if d.error { // ran out of data
		return nil, badData
	}

	var extend string
	rest := d.rest()
	if len(rest) > 2 && rest[0] == '\n' && rest[len(rest)-1] == '\n' {
		extend = string(rest[1 : len(rest)-1])
	}

	// Now we can build up a useful data structure.
	// First the zone information.
	//  utcoff[4] isdst[1] nameindex[1]
	nzone := cnt[NZone]
	if nzone == 0 {
		// Reject tzdata files with no zones. There's nothing useful in them.
		// This also avoids a panic later when we add and then use a fake transition (golang.org/issue/29437).
		return nil, badData
	}
	eras := make([]Era, nzone)
	for i := range eras {
		var ok bool
		var n uint32
		if n, ok = zonedata.big4(); !ok {
			return nil, badData
		}
		if uint32(int(n)) != n {
			return nil, badData
		}
		eras[i].Offset = int(int32(n))
		var b byte
		if b, ok = zonedata.byte(); !ok {
			return nil, badData
		}
		eras[i].IsDST = b != 0
		if b, ok = zonedata.byte(); !ok || int(b) >= len(abbrev) {
			return nil, badData
		}
		eras[i].Name = byteString(abbrev[b:])
		if runtime.GOOS == "aix" && len(name) > 8 && (name[:8] == "Etc/GMT+" || name[:8] == "Etc/GMT-") {
			// There is a bug with AIX 7.2 TL 0 with files in Etc,
			// GMT+1 will return GMT-1 instead of GMT+1 or -01.
			if name != "Etc/GMT+0" {
				// GMT+0 is OK
				eras[i].Name = name[4:]
			}
		}
	}

	// Get Unix timestamps for the edges of a two-year period,
	// centered around current date.
	// These will be used to "filter" recent transitions.
	now := time.Now()
	past := now.AddDate(-1, 0, 0)
	pastEpoch := past.Unix()
	future := now.AddDate(1, 0, 0)
	futureEpoch := future.Unix()

	// Now the transition time info.
	usefullTrans := 0
	tx := make([]EraTrans, cnt[NTime])
	for i := range tx {
		// Get time of transition (in Unix format)
		var n int64
		if !is64 {
			if n4, ok := txtimes.big4(); !ok {
				return nil, badData
			} else {
				n = int64(int32(n4))
			}
		} else {
			if n8, ok := txtimes.big8(); !ok {
				return nil, badData
			} else {
				n = int64(n8)
			}
		}

		// Skip transition if in past or far into the future.
		if n < pastEpoch || n > futureEpoch {
			continue
		}

		tx[usefullTrans].When = n

		if int(txzones[i]) >= len(eras) {
			return nil, badData
		}
		tx[usefullTrans].Index = txzones[i]

		if i < len(isstd) {
			tx[usefullTrans].Isstd = (isstd[i] != 0)
		}

		if i < len(isutc) {
			tx[usefullTrans].Isutc = (isutc[i] != 0)
		}

		usefullTrans++
	}

	l := &TZdata{Eras: eras, Trans: tx[:usefullTrans], Name: name, Extend: extend}

	return l, nil
}
