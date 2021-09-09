// +build linux

package machineid

import (
	"os"
	"path"
	"regexp"
	"errors"
	"strings"
)

const (
	// dbusPath is the default path for dbus machine id.
	dbusPath = "/var/lib/dbus/machine-id"
	// dbusPathEtc is the default path for dbus machine id located in /etc.
	// Some systems (like Fedora 20) only know this path.
	// Sometimes it's the other way round.
	dbusPathEtc = "/etc/machine-id"

	mountInfoPath = "/proc/self/mountinfo"

	ctrlGroupPath = "/proc/self/cgroup"

	dockerEnvFile = "/.dockerenv"

	lengthOfContainerId = 64
)

var LxcCtrlGroupKeywords = []string{
	"docker",
	"kubepods",
}

func isContainer() bool {
	_, err := os.Stat(dockerEnvFile)
	if ! os.IsNotExist(err) {
		return true
	}

	bytes, err := readFile(ctrlGroupPath)
	if err != nil {
		return false
	}

	info := string(bytes)
	for _, keyword := range LxcCtrlGroupKeywords {
		if strings.Contains(info, keyword) {
			return true
		}
	}

	return false
}

func extractMountInfoId(mountInfo string) string {
	lines := strings.Split(mountInfo, "\n")
	pattern := regexp.MustCompile("upperdir=(.+?),")

	for _, line := range lines {
		upperDir := pattern.FindString(line)
		if upperDir == "" {
			continue
		}

		dir, _ := path.Split(upperDir)
		return path.Base(dir)
	}

	return ""
}

func extractCtrlGroupId(ctrlGroup string) string {
	lines := strings.Split(ctrlGroup, "\n")

	basePattern := regexp.MustCompile("/([a-f0-9]{64})$")
	scopePattern := regexp.MustCompile("/.+-(.+?)\\.scope$")

	for _, line := range lines {
		// See https://stackoverflow.com/questions/20010199/how-to-determine-if-a-process-runs-inside-lxc-docker
		// systemd Docker
		// 1:cpuset:/kubepods.slice/kubepods-burstable.slice/kubepods-burstable-podf8e932ad_5487_4514_a5be_b75ad1b7a6ce.slice/crio-ee55f03bd921c55955d8995a0adbb9f19352603a637ea27f6ca8397b715435eb.scope
		// 5:net_cls:/system.slice/docker-afd862d2ed48ef5dc0ce8f1863e4475894e331098c9a512789233ca9ca06fc62.scope
		slices := scopePattern.FindStringSubmatch(line)
		if len(slices) >= 2 && len(slices[1]) == lengthOfContainerId {
			return slices[1]
		}

		// Non-systemd Docker
		// 5:net_prio,net_cls:/docker/de630f22746b9c06c412858f26ca286c6cdfed086d3b302998aa403d9dcedc42
		// 3:net_cls:/kubepods/burstable/pod5f399c1a-f9fc-11e8-bf65-246e9659ebfc/9170559b8aadd07d99978d9460cf8d1c71552f3c64fefc7e9906ab3fb7e18f69
		slices = basePattern.FindStringSubmatch(line)
		if len(slices) >= 2 && len(slices[1]) == lengthOfContainerId {
			return slices[1]
		}
	}

	return ""
}

func readMountInfo() (string, error) {
	lines, err := readFile(mountInfoPath)
	if err != nil {
		return "", err
	}

	id := extractMountInfoId(string(lines))
	if id != "" {
		return id, nil
	}

	return "", errors.New("no storage driver found in mount table")
}

func readCtrlGroups() (string, error) {
	lines, err := readFile(ctrlGroupPath)
	if err != nil {
		return "", err
	}

	id := extractCtrlGroupId(string(lines))
	if id != "" {
		return id, nil
	}

	return "", errors.New("cannot find any id in cgroup file")
}

func getContainerId() (string, error) {
	id, err := readMountInfo()
	if err != nil {
		// try fallback path
		id, err = readCtrlGroups()
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func getHostId() (string, error) {
	id, err := readFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = readFile(dbusPathEtc)
	}
	if err != nil {
		return "", err
	}
	return trim(string(id)), nil
}

// machineID returns the uuid specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// If there is an error reading the files an empty string is returned.
// See https://unix.stackexchange.com/questions/144812/generate-consistent-machine-unique-id
func machineID() (string, error) {
	id, err := getHostId()
	if err != nil && isContainer() {
		id, err = getContainerId()
	}
	if err != nil {
		return "", err
	}
	return id, nil
}
