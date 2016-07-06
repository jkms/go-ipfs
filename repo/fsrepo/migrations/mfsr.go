package mfsr

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const VersionFile = "version"

type RepoPath string

func (rp RepoPath) VersionFile() string {
	return path.Join(string(rp), VersionFile)
}

func (rp RepoPath) Version() (int, error) {
	if rp == "" {
		return 0, fmt.Errorf("invalid repo path \"%s\"", rp)
	}

	fn := rp.VersionFile()
	if _, err := os.Stat(fn); os.IsNotExist(err) {
		return 0, VersionFileNotFound(rp)
	}

	c, err := ioutil.ReadFile(fn)
	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(c))
	return strconv.Atoi(s)
}

func (rp RepoPath) CheckVersion(version int) error {
	v, err := rp.Version()
	if err != nil {
		return err
	}

	if v != version {
		return fmt.Errorf("versions differ (expected: %s, actual:%s)", version, v)
	}

	return nil
}

func (rp RepoPath) WriteVersion(version int) error {
	fn := rp.VersionFile()
	return ioutil.WriteFile(fn, []byte(fmt.Sprintf("%d\n", version)), 0644)
}

type VersionFileNotFound string

func (v VersionFileNotFound) Error() string {
	return "no version file in repo at " + string(v)
}

func TryMigrating(tovers int) error {
	if !YesNoPrompt("Run migrations automatically? [y/N]") {
		return fmt.Errorf("please run the migrations manually")
	}

	return RunMigration(tovers)
}

func YesNoPrompt(prompt string) bool {
	var s string
	for i := 0; i < 3; i++ {
		fmt.Printf("%s ", prompt)
		fmt.Scanf("%s", &s)
		switch s {
		case "y", "Y":
			return true
		case "n", "N":
			return false
		case "":
			return false
		}
		fmt.Println("Please press either 'y' or 'n'")
	}

	return false
}
