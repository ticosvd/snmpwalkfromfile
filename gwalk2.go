/*
Version: 0.0.1



*/

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gosnmp/gosnmp"
)

// Struct for snmp walking from file
type OIDSFile struct {
	file      string
	oids      []string
	target    string
	community string
}

// Func creates slice of OID from snmp file
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

// Func prints result of snmpwaling by oid , it gets from gosnmp example
func (o OIDSFile) printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	return nil
}

// Func gets result from SNMPwalk by oid, main parts from gosnmp
func (o OIDSFile) SNMPWalker(oid string) error {

	gosnmp.Default.Target = o.target
	gosnmp.Default.Community = o.community
	gosnmp.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
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
	for {
		for index, oid := range o.oids {
			log.Println("Current Index", index, "OID", oid)
			err := o.SNMPWalker(oid)
			if err != nil {
				log.Printf("Error in OID", oid, "and error", err)
				return err
			}
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return nil
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s [-community=<community>] host \n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community, snmpfile string
	var interval int
	flag.StringVar(&community, "community", "public", "the community string for device")
	flag.IntVar(&interval, "interval", 60, "Interval in seconds")
	flag.StringVar(&snmpfile, "file", "snmp2.mibs", "File in Mibs")

	flag.Parse()

	log.Println(snmpfile)
	oidsFromFile := OIDSFile{file: snmpfile}
	toids, err := oidsFromFile.GettingOIDS()
	log.Println(oidsFromFile.oids)
	if err != nil {
		log.Fatal(err)
	}

	oidsFromFile.oids = make([]string, len(toids))
	co := copy(oidsFromFile.oids, toids)
	if co < 1 {
		log.Fatal(errors.New("No slice with oids !!!"))
	}
	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	oidsFromFile.target = flag.Args()[0]
	oidsFromFile.community = community

	err = oidsFromFile.StartSNMPWalker(interval)
	if err != nil {
		log.Fatal(err)
	}

}
