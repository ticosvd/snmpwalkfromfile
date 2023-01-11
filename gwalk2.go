/*
Version: 0.0.2
Devide code to two files:
gwalk2.go - main and http server for snmpwalk
snmppidsfile.go - snmpwalker , snmpbulkget
/dbmap/dbmap.go - structs  for tdb


*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func httpserver(w http.ResponseWriter, _ *http.Request) {
	dd := d
	namex := make([]string, 0, len(dd.Oiddata))
	namey := make([]opts.LineData, 0, len(dd.Oiddata))
	for k, v := range dd.Oiddata {
		namex = append(namex, k.Format("2006-1-2 15:4:5"))
		namey = append(namey, opts.LineData{Value: fmt.Sprintf("%v", v.Value)})

	}

	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title:    "MP CPU Loading",
			Subtitle: "Line chart rendered by the http server this time",
		}))

	// Put data into instance
	line.SetXAxis(namex).
		AddSeries("SNMP oid ", namey).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)

}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("Version 0.0.1\n")
		fmt.Printf("   %s [-community=<community>] host \n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community, snmpfile string
	var interval int
	var walk string
	flag.StringVar(&community, "community", "public", "the community string for device")
	flag.IntVar(&interval, "interval", 0, "Interval in seconds")
	flag.StringVar(&snmpfile, "file", "snmp2.mibs", "File in Mibs")
	flag.StringVar(&walk, "walk", "yes", "Choose Walk or Bulk")

	flag.Parse()

	log.Println(snmpfile)
	oidsFromFile := OIDSFile{file: snmpfile}
	toids, err := oidsFromFile.GettingOIDS()
	log.Println(oidsFromFile.oids)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", httpserver)
	go http.ListenAndServe(":8081", nil)

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
	if walk == "yes" {
		err = oidsFromFile.StartSNMPWalker(interval)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = oidsFromFile.GetSNMPBulk(interval)

		if err != nil {
			log.Fatal(err)
		}
	}
}
