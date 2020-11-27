package tickscript_test

import "testing"
import "csv"
import "contrib/bonitoo-io/tickscript"
import "influxdata/influxdb/monitor"
import "influxdata/influxdb/schema"

option now = () => (2020-11-25T14:05:30Z)

option monitor.write = (tables=<-) => tables
option monitor.log = (tables=<-) => tables

statusesData = "
#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,double,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,0,2020-11-25T14:04:25.756624987Z,1.53716449678207,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka01,ft
,,1,2020-11-25T14:04:25.856624989Z,8.33716449678206,rate-check,Rate Check,KafkaMsgRate,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,2,2020-11-25T14:04:25.756624987Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert: ok - 1.53716449678207,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka01,ft
,,3,2020-11-25T14:04:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,rate-check,Rate Check,_message,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,long,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,4,2020-11-25T14:04:25.756624987Z,1606313047768317098,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka01,ft
,,5,2020-11-25T14:04:25.856624989Z,1606313047768317097,rate-check,Rate Check,_source_timestamp,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,6,2020-11-25T14:04:25.756624987Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka01,ft
,,7,2020-11-25T14:04:25.856624989Z,some detail: myrealm=ft,rate-check,Rate Check,details,warn,statuses,testm,custom,kafka07,ft
,,8,2020-11-25T14:04:25.756624987Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka01,ft
,,9,2020-11-25T14:04:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,warn,statuses,testm,custom,kafka07,ft
"

// override source for existing statuses
option tickscript._sssource = () =>
  csv.from(csv: statusesData)

inData = "
#group,false,false,false,false,true,true,true,true
#datatype,string,long,dateTime:RFC3339,double,string,string,string,string
#default,_result,,,,,,,
,result,table,_time,_value,_field,_measurement,host,realm
,,0,2020-11-25T14:05:03.477635916Z,1.819231109049999,kafka_message_in_rate,testm,kafka07,ft
,,0,2020-11-25T14:05:04.541635074Z,1.635878190200181,kafka_message_in_rate,testm,kafka07,ft
,,0,2020-11-25T14:05:05.623191313Z,39.33716449678206,kafka_message_in_rate,testm,kafka07,ft
,,0,2020-11-25T14:05:06.696061106Z,26.33716449678206,kafka_message_in_rate,testm,kafka07,ft
,,0,2020-11-25T14:05:07.768317097Z,8.33716449678206,kafka_message_in_rate,testm,kafka07,ft
,,0,2020-11-25T14:05:08.868317091Z,1.33716449678206,kafka_message_in_rate,testm,kafka07,ft
"

outData = "
#group,false,false,false,true,true,true,true,false,true,false,true,false,true,false,true
#datatype,string,long,double,string,string,string,string,string,string,long,string,string,string,string,string
#default,_result,,,,,,,,,,,,,,
,result,table,KafkaMsgRate,_check_id,_check_name,_level,_measurement,_message,_source_measurement,_source_timestamp,_type,details,host,id,realm
,,0,1.819231109049999,rate-check,Rate Check,ok,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.819231109049999,testm,1606313103477635916,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft
,,0,1.33716449678206,rate-check,Rate Check,ok,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.33716449678206,testm,1606313108868317091,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft
,,1,39.33716449678206,rate-check,Rate Check,crit,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 39.33716449678206,testm,1606313105623191313,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft
,,1,26.33716449678206,rate-check,Rate Check,crit,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 26.33716449678206,testm,1606313106696061106,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft
,,2,8.33716449678206,rate-check,Rate Check,warn,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,testm,1606313107768317097,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft
"

check = {
  _check_id: "rate-check",
  _check_name: "Rate Check",
  _type: "custom", // tickscript?
  tags: {},
}

metric_type = "kafka_message_in_rate"
tier = "ft"
h_threshold = 10
w_threshold = 5
l_threshold = .002

tickscript_alert = (table=<-) => table
	|> range(start: 2020-11-25T14:05:00Z)
    |> filter(fn: (r) => r._field == metric_type and r.realm == tier)
    |> group(columns: ["_measurement", "host", "realm"])
    |> schema.fieldsAsCols()
    |> drop(columns: ["_start", "_stop"])
    |> tickscript.as(column: metric_type, as: "KafkaMsgRate")
    |> tickscript.alert(
        check: check,
        id: (r) => "Realm: ${r.realm} - Hostname: ${r.host} / Metric: ${metric_type} threshold alert",
        message: (r) => "${r.id}: ${r._level} - ${string(v:r.KafkaMsgRate)}",
        details: (r) => "some detail: myrealm=${r.realm}",
        crit: (r) => r.KafkaMsgRate > h_threshold or r.KafkaMsgRate < l_threshold,
        warn: (r) => r.KafkaMsgRate > w_threshold or r.KafkaMsgRate < l_threshold
    )
    |> drop(columns: ["_time"])

test _tickscript_alert = () => ({
	input: testing.loadStorage(csv: inData),
	want: testing.loadMem(csv: outData),
	fn: tickscript_alert,
})
