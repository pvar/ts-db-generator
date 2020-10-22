package tzdata

import (
	"errors"
	"syscall"
)

// maxFileSize is the max permitted size of files read by readFile.
// As reference, the zoneinfo.zip distributed by Go is ~350 KB,
// so 10MB is overkill.
const maxFileSize = 10 << 20

// loadLocation returns the Location with the given name from one of
// the specified sources. See loadTzinfo for a list of supported sources.
// The first timezone data matching the given name that is successfully loaded
// and parsed is returned as a Location.
func readTZfile(name string) (z *TZdata, firstErr error) {
	var file string

	for _, source := range []string{"/usr/share/zoneinfo/", "/usr/share/lib/zoneinfo/"} {
		file = source + "/" + name

		var rawTZdata, err = loadFile(file)

		if err == nil {
			data, err := parseRawTZdata(name, rawTZdata)
			if err == nil {
				return data, nil
			}
			return nil, err
		}
	}
	return nil, errors.New("tzdata: unknown time zone " + name)
}

// loadFile reads and returns the content of the named file.
// It is a trivial implementation of ioutil.ReadFile,
// reimplemented to avoid depending on io/ioutil or os.
// It returns an error if name exceeds maxFileSize bytes.
func loadFile(name string) ([]byte, error) {
	f, err := open(name)
	if err != nil {
		return nil, err
	}
	defer closefd(f)

	var (
		buf [4096]byte
		ret []byte
		n   int
	)

	for {
		n, err = read(f, buf[:])
		if n > 0 {
			ret = append(ret, buf[:n]...)
		}
		if n == 0 || err != nil {
			break
		}
		if len(ret) > maxFileSize {
			return nil, errors.New("tzdata: timezone file too big")
		}
	}
	return ret, err
}

func open(name string) (uintptr, error) {
	fd, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		return 0, err
	}
	return uintptr(fd), nil
}

func read(fd uintptr, buf []byte) (int, error) {
	return syscall.Read(int(fd), buf)
}

func closefd(fd uintptr) {
	syscall.Close(int(fd))
}
