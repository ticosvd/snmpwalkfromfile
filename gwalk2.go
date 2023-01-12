/*
Version: 0.0.3
Devide code to two files:
gwalk2.go - main and http server for snmpwalk
snmppidsfile.go - snmpwalker , snmpbulkget
/dbmap/dbmap.go - structs  for tdb
/logger/logger.go - for logging messages

*/

package main

import (
	"errors"
	"flag"
	"fmt"
	"gwalk2byoids/logger"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

var loglevel logger.Logstruct

func httpserver(w http.ResponseWriter, _ *http.Request) {
	dd := d
	namext := make([]time.Time, 0, len(dd.Oiddata))
	namex := make([]string, 0, len(dd.Oiddata))
	namey := make([]opts.LineData, 0, len(dd.Oiddata))
	//namey := make([]opts.KlineData, 0, len(dd.Oiddata))

	for k := range dd.Oiddata {
		namext = append(namext, k)

	}

	loglevel.Logger("Before", namext)
	sort.Slice(namext, func(i, j int) bool { return namext[i].Before(namext[j]) })

	loglevel.Logger("After", namext)
	for _, v := range namext {

		//	namex = append(namex, namext[i].Format("2006-1-2 15:4:5"))
		namex = append(namex, v.Format("2006-01-02 15:04:05"))
		namey = append(namey, opts.LineData{Value: fmt.Sprintf("%v", dd.Oiddata[v].Value)})
		//namey = append(namey, opts.KlineData{Value: fmt.Sprintf("%v", dd.Oiddata[v].Value)})
	}

	line := charts.NewLine()

	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeWesteros}),
		charts.WithTitleOpts(opts.Title{
			Title: "MP CPU Loading",
			//	Subtitle: "SSS" + namext[0].Format("2006-1-2 15:4:5") + "BBB",
		}),

		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: true,
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "slider",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	// Put data into instance
	line.SetXAxis(namex).
		AddSeries("lines", namey)
		//AddSeries("SNMP oid ", namey).
		//		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: true}))
	line.Render(w)

}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("Version 0.0.2\n")
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
		//	loglevel.Loglevel = "DEBUG"
		http.HandleFunc("/", httpserver)
		go http.ListenAndServe(":8081", nil)
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
