/*
Version 0.0.2
Includes struct OIDSFile and functions for the struct
*/
package main

import (
	"bufio"
	"fmt"
	"gwalk2byoids/dbmap"
	"log"
	"os"
	"time"

	"github.com/gosnmp/gosnmp"
)

// the variable collect tdb data
var d dbmap.Databytime

// Struct for snmp walking from file
type OIDSFile struct {
	file      string
	oids      []string
	target    string
	community string
}

// Func gets all oids from the file and append ones to toids
func (o OIDSFile) GettingOIDS() ([]string, error) {
	var toids []string
	f, err := os.Open(o.file)

	if err != nil {
		log.Println("The problem with openning file", err)
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		//		log.Println("OIDS", scanner.Text())
		toids = append(toids, scanner.Text())
	}

	return toids, nil

}

// Func prints result of snmpwaling by oid , it gets from gosnmp example and set dbmapstruct
func (o OIDSFile) printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	d.Add(time.Now(), pdu.Name, pdu.Value)
	return nil
}

// Func gets result from SNMPwalk by oid, main parts from gosnmp
func (o OIDSFile) SNMPWalker(oid string) error {

	gosnmp.Default.Target = o.target
	gosnmp.Default.Community = o.community
	gosnmp.Default.Timeout = time.Duration(5 * time.Second) // Timeout better suited to walking
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		return err
	}
	defer gosnmp.Default.Conn.Close()

	err = gosnmp.Default.BulkWalk(oid, o.printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		return err
	}
	return nil
}

// Func strats inifinite loops and snmpwpooling oids from the slice
func (o OIDSFile) StartSNMPWalker(interval int) error {
	iter := 1
	d.Oiddata = make(map[time.Time]dbmap.Valueoid)
	for {
		log.Println("Current iteration:", iter)
		for index, oid := range o.oids {
			log.Println("Current Index", index, "OID", oid)
			err := o.SNMPWalker(oid)
			if err != nil {
				log.Printf("Error in OID", oid, "and error", err)
				return err
			}
		}
		iter = iter + 1
		fmt.Println("------------------------------------------------")
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return nil
}

// Func starts GetBulk SNMP requests
func (o OIDSFile) GetSNMPB(oids []string) error {

	gosnmp.Default.Target = o.target
	gosnmp.Default.Community = o.community
	gosnmp.Default.Timeout = time.Duration(5 * time.Second) // Timeout better suited to walking
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		return err
	}
	defer gosnmp.Default.Conn.Close()

	//snmpresult, err := gosnmp.Default.GetBulk(oids, 1, 1)
	_, err = gosnmp.Default.GetBulk(oids, 1, 1)
	if err != nil {
		fmt.Printf("Bulk Error: %v\n", err)
		return err
	}
	// log.Println("Current Bulk get", snmpresult)

	return nil
}

// Func starts GetBulk SNMP requests
func (o OIDSFile) GetSNMPBulk(interval int) error {
	iter := 1
	for {
		log.Println("Current iteration:", iter)

		for i := 0; i < len(o.oids); i += 60 {
			fmt.Println("curr index", i)
			b := i + 59
			if b < len(o.oids)-1 {

				err := o.GetSNMPB(o.oids[i:b])
				if err != nil {
					fmt.Printf("SNMP GET Error: %v\n", err)
					return err
				}

			} else {

				err := o.GetSNMPB(o.oids[i:len(o.oids)])
				if err != nil {
					fmt.Printf("Bulk Error: %v\n", err)
					return err
				}

			}

		}

		iter = iter + 1
		fmt.Println("------------------------------------------------")
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return nil
}
