/*
Version 0.0.1
It has sturcts :
Valueoid
Databytime and fucntions for the struct
*/
package dbmap

import (
	"log"
	"time"
)

type Valueoid struct {
	Oid   string
	Value interface{}
}

// the main struct which includes data : map[timetype]:{Oid:xxx, Value:yyy }
type Databytime struct {
	Oiddata map[time.Time]Valueoid
}

func (d Databytime) Add(t time.Time, oid string, value interface{}) error {
	//	if d.Oiddata == nil {
	//d.Oiddata = make(map[time.Time]Valueoid)
	//	}
	valoid := Valueoid{Oid: oid, Value: value}
	d.Oiddata[time.Now()] = valoid

	return nil
}

func (d Databytime) PrintAll() {
	for k, v := range d.Oiddata {
		log.Printf("KEY is %s and value: OID : %s, VALUE : %v \n", k.Format("2006-01-02 03:04:05"), v.Oid, v.Value)
	}
}
