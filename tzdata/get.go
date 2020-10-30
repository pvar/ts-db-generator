package tzdata

import (
    "os"
    "bufio"
    "strings"
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

func GetList () (version string, timezones map[string]string, err error) {
    file, err := os.Open(source_path + "tzdata.zi")
    if err != nil {
        return "", nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    // get version string from first line
    scanner.Scan()
    substr := strings.Fields(scanner.Text())
    version = substr[2]

    // scan the rest of the lines
    timezones = make(map[string]string, 500)
    for scanner.Scan() {
        substr = strings.Fields(scanner.Text())
        if substr[0] == "Z" {
            timezones[substr[1]] = substr[1]
        }

        if substr[0] == "L" {
            timezones[substr[2]] = substr[1]
        }
    }

    // check if scanning stopped due to some error
    if err = scanner.Err(); err != nil {
        return "", nil, err
    }

    return version, timezones, nil
}
