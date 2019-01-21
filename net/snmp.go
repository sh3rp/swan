package net

import (
	"fmt"
	"time"

	"github.com/sh3rp/swan/version"
	"github.com/soniah/gosnmp"
)

var OIDS = map[string]string{
	"sysDescr":  "1.3.6.1.2.1.1.1.0",
	"sysUptime": "1.3.6.1.2.1.1.3.0",
	"sysName":   "1.3.6.1.2.1.1.5.0",
	"ifIndex":   "1.3.6.1.2.1.2.2.1.1",
	"ifDesc":    "1.3.6.1.2.1.2.2.1.2",
	"ifAlias":   "1.3.6.1.2.1.31.1.1.1.18",
	// juniper specific
	"ifIn1SecRate":   "1.3.6.1.4.1.2636.3.3.1.1.1",
	"ifOut1SecRate":  "1.3.6.1.4.1.2636.3.3.1.1.4",
	"ifJnxInErrors":  "1.3.6.1.4.1.2636.3.3.1.1.9",
	"ifJnxOutErrors": "1.3.6.1.4.1.2636.3.3.1.1.24",
}

type SwitchManager interface {
	GetVersion() (SwitchInfo, error)
	GetIfs() ([]SwitchInterface, error)
	GetIfStats(SwitchInterface) (SwitchIfStats, error)
}

type SwitchInfo struct {
	Hostname  string
	OSVersion version.OSVersion
}

type SwitchInterface struct {
	Name      string
	Label     string
	SnmpIndex int
	Status    string
}

type SwitchIfStats struct {
	IfBitsInPerSecond  uint
	IfBitsOutPerSecond uint
	IfPacketsIn        uint
	IfPacketsOut       uint
	IfInErrors         uint64
	IfOutErrors        uint64
}

type switchManager struct {
	snmp gosnmp.GoSNMP
}

func NewSwitchManager(ip string, community string) SwitchManager {
	return switchManager{gosnmp.GoSNMP{
		Target:    ip,
		Port:      uint16(161),
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
	}}
}

func (sm switchManager) GetVersion() (SwitchInfo, error) {
	err := sm.snmp.Connect()

	if err != nil {
		return SwitchInfo{}, err
	}

	defer sm.snmp.Conn.Close()

	oids := []string{OIDS["sysDescr"], OIDS["sysName"]}

	result, err := sm.snmp.Get(oids)

	if err != nil {
		return SwitchInfo{}, err
	}

	return SwitchInfo{
		Hostname:  string(result.Variables[1].Value.([]byte)),
		OSVersion: version.GetVersion(string(result.Variables[0].Value.([]byte))),
	}, nil
}

func (sm switchManager) GetIfs() ([]SwitchInterface, error) {
	err := sm.snmp.Connect()
	if err != nil {
		return nil, err
	}

	defer sm.snmp.Conn.Close()

	var ifs []int

	sm.snmp.BulkWalk(OIDS["ifIndex"], func(pdu gosnmp.SnmpPDU) error {
		ifs = append(ifs, pdu.Value.(int))
		return nil
	})

	var swIf []SwitchInterface

	for _, intf := range ifs {
		oids := []string{
			getOid("ifDesc", intf),
			getOid("ifAlias", intf),
		}
		results, _ := sm.snmp.Get(oids)
		switchInterface := SwitchInterface{
			Name:      string(results.Variables[0].Value.([]byte)),
			Label:     string(results.Variables[1].Value.([]byte)),
			SnmpIndex: intf,
		}
		swIf = append(swIf, switchInterface)
	}

	return swIf, nil
}

func (sm switchManager) GetIfStats(swIf SwitchInterface) (SwitchIfStats, error) {
	err := sm.snmp.Connect()
	if err != nil {
		return SwitchIfStats{}, err
	}

	defer sm.snmp.Conn.Close()

	intf := swIf.SnmpIndex
	oids := []string{
		getOid("ifIn1SecRate", intf),
		getOid("ifOut1SecRate", intf),
		getOid("ifJnxInErrors", intf),
		getOid("ifJnxOutErrors", intf),
	}

	results, err := sm.snmp.Get(oids)

	return SwitchIfStats{
		IfBitsInPerSecond:  results.Variables[0].Value.(uint),
		IfBitsOutPerSecond: results.Variables[1].Value.(uint),
		IfInErrors:         results.Variables[2].Value.(uint64),
		IfOutErrors:        results.Variables[3].Value.(uint64),
	}, nil
}

func getOid(label string, idx int) string {
	return fmt.Sprintf("%s.%d", OIDS[label], idx)
}
