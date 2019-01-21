package version

import (
	"regexp"
)

type OSType int

const (
	UNKNOWN = iota
	CISCO_IOS
	CISCO_NXOS
	JUNOS
	ARISTA
)

func (ost OSType) String() string {
	switch ost {
	case CISCO_IOS:
		return "Cisco IOS"
	case CISCO_NXOS:
		return "Cisco NX-OS"
	case JUNOS:
		return "JUNOS"
	case ARISTA:
		return "Arista"
	}
	return "Unknown"
}

type OSVersion struct {
	OSType  OSType
	Version string
}

var ciscoIOSSignature = regexp.MustCompile("Cisco Internetworking")
var ciscoNXOSSignature = regexp.MustCompile("NX-OS")
var junOSSignature = regexp.MustCompile("Juniper")
var junOSSVersion = regexp.MustCompile("JUNOS\\s([0-9.a-zA-Z]+)\\,")

func GetVersion(str string) OSVersion {
	osType := OSType(UNKNOWN)
	version := ""

	if ciscoIOSSignature.MatchString(str) {
		osType = CISCO_IOS
		version = "unknown"
	} else if ciscoNXOSSignature.MatchString(str) {
		osType = CISCO_NXOS
		version = "unknown"
	} else if junOSSignature.MatchString(str) {
		osType = JUNOS
		v := junOSSVersion.FindStringSubmatch(str)
		version = v[1]
	}

	return OSVersion{OSType: osType, Version: version}
}
