# snmpwalkfromfile

Version 0.0.2
   gwalk2 [-community=<community>] host 
     host      - the host to walk/scan
     oid       - the MIB/Oid defining a subtree of values

  -community string
    	the community string for device (default "public")
  -file string
    	File in Mibs (default "snmp2.mibs")
  -interval int
    	Interval in seconds
  -walk string
    	Choose Walk or Bulk (default "yes")

Bulk mode is only permanently getting data (emulation of Network Management System), It doesn't collect data.


