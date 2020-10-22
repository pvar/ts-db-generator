package tzdata

import (
	"errors"
)

func GetData(location string) (*TZdata, error) {
	if len(location) == 0 {
		return nil, errors.New("tzdata: empty location name")
	}

	if containsDotDot(location) || location[0] == '/' || location[0] == '\\' {
		// No valid IANA Time Zone name contains a single dot,
		// much less dot dot. Likewise, none begin with a slash.
		return nil, errors.New("tzdata: invalid location name")
	}

	data, err := readTZfile(location)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func containsDotDot(s string) bool {
	if len(s) < 2 {
		return false
	}
	for i := 0; i < len(s)-1; i++ {
		if s[i] == '.' && s[i+1] == '.' {
			return true
		}
	}
	return false
}
