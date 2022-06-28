package fileutils

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ReadLine reads certain line from the file.
func ReadLine(file *os.File, lineNum int) (string, error) {
	var line int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line++

		if line == lineNum {
			return strings.TrimSpace(scanner.Text()), scanner.Err()
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", io.EOF
}

// CountLines counts lines from file.
func CountLines(r io.Reader) (int, error) {
	var line int
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return line, nil
}

// ApplicationDir returns best base directory for specific OS.
func ApplicationDir(subdir ...string) string {
	for i := range subdir {
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			subdir[i] = cases.Title(language.English).String(subdir[i])
		} else {
			subdir[i] = cases.Lower(language.English).String(subdir[i])
		}
	}

	var appdir string

	home := os.Getenv("HOME")

	switch runtime.GOOS {
	case "windows":
		// Windows standards: https://msdn.microsoft.com/en-us/library/windows/apps/hh465094.aspx?f=255&MSPPError=-2147217396
		for _, env := range []string{"AppData", "AppDataLocal", "UserProfile", "Home"} {
			val := os.Getenv(env)
			if val != "" {
				appdir = val
				break
			}
		}
	case "darwin":
		// Mac standards: https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/MacOSXDirectories/MacOSXDirectories.html
		appdir = filepath.Join(home, "Library", "Application Support")
	case "linux":
		fallthrough
	default:
		// Linux standards: https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html
		appdir = os.Getenv("XDG_DATA_HOME")
		if appdir == "" && home != "" {
			appdir = filepath.Join(home, ".local", "share")
		}
	}

	return filepath.Join(append([]string{appdir}, subdir...)...)
}

// IsFileExist checks if file with given name exists in path.
func IsFileExist(path, fName string) (bool, error) {
	name := path
	if name[len(name)-1:] != "/" {
		name += "/"
	}

	name += fName
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// CreateFile creates a new file in path.
func CreateFile(path, fName string) error {
	name := path
	if name[len(name)-1:] != "/" {
		name += "/"
	}

	name += fName
	f, err := os.Create(name)
	if err != nil {
		return err
	}

	_ = f.Close()
	return nil
}
