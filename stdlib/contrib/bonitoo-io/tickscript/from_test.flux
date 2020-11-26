package tickscript_test

import "testing"
import "csv"
import "contrib/bonitoo-io/tickscript"
import "experimental"
import "experimental/array"

option now = () => (2020-11-25T14:06:00Z)

// TODO remove
topicsData_notKeyed = "
#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,double,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,0,2020-11-25T14:05:25.856485929Z,39.33716449678206,rate-check,Rate Check,KafkaMsgRate,crit,statuses,testm,custom,kafka07,ft
,,0,2020-11-25T14:05:25.856575405Z,26.33716449678206,rate-check,Rate Check,KafkaMsgRate,crit,statuses,testm,custom,kafka07,ft
,,1,2020-11-25T14:05:25.576319539Z,1.419231109050123,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka01,ft
,,1,2020-11-25T14:05:25.856319359Z,1.819231109049999,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka01,ft
,,1,2020-11-25T14:05:25.856444173Z,1.635878190200181,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka01,ft
,,2,2020-11-25T14:05:25.856624989Z,8.33716449678206,rate-check,Rate Check,KafkaMsgRate,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,3,2020-11-25T14:05:25.856485929Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 39.33716449678206,rate-check,Rate Check,_message,crit,statuses,testm,custom,kafka07,ft
,,3,2020-11-25T14:05:25.856575405Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 26.33716449678206,rate-check,Rate Check,_message,crit,statuses,testm,custom,kafka07,ft
,,4,2020-11-25T14:05:25.576319539Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert: ok - 1.419231109050123,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka01,ft
,,4,2020-11-25T14:05:25.856319359Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.819231109049999,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka01,ft
,,4,2020-11-25T14:05:25.856444173Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.635878190200181,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka01,ft
,,5,2020-11-25T14:05:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,rate-check,Rate Check,_message,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,long,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,6,2020-11-25T14:05:25.856485929Z,1606313105623191313,rate-check,Rate Check,_source_timestamp,crit,statuses,testm,custom,kafka07,ft
,,6,2020-11-25T14:05:25.856575405Z,1606313106696061106,rate-check,Rate Check,_source_timestamp,crit,statuses,testm,custom,kafka07,ft
,,7,2020-11-25T14:05:25.576319539Z,1606313101487639160,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka01,ft
,,7,2020-11-25T14:05:25.856319359Z,1606313103477635916,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka01,ft
,,7,2020-11-25T14:05:25.856444173Z,1606313104541635074,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka01,ft
,,8,2020-11-25T14:05:25.856624989Z,1606313107768317097,rate-check,Rate Check,_source_timestamp,warn,statuses,testm,custom,kafka07,ft

#group,false,false,false,false,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm
,,9,2020-11-25T14:05:25.856485929Z,some detail: myrealm=ft,rate-check,Rate Check,details,crit,statuses,testm,custom,kafka07,ft
,,9,2020-11-25T14:05:25.856575405Z,some detail: myrealm=ft,rate-check,Rate Check,details,crit,statuses,testm,custom,kafka07,ft
,,10,2020-11-25T14:05:25.576319539Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka01,ft
,,10,2020-11-25T14:05:25.856319359Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka01,ft
,,10,2020-11-25T14:05:25.856444173Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka01,ft
,,11,2020-11-25T14:05:25.856624989Z,some detail: myrealm=ft,rate-check,Rate Check,details,warn,statuses,testm,custom,kafka07,ft
,,12,2020-11-25T14:05:25.856485929Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,crit,statuses,testm,custom,kafka07,ft
,,12,2020-11-25T14:05:25.856575405Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,crit,statuses,testm,custom,kafka07,ft
,,13,2020-11-25T14:05:25.576319539Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka01,ft
,,13,2020-11-25T14:05:25.856319359Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka01,ft
,,13,2020-11-25T14:05:25.856444173Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka01,ft
,,14,2020-11-25T14:05:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,warn,statuses,testm,custom,kafka07,ft
,,15,2020-11-25T14:05:25.856485929Z,TESTING,rate-check,Rate Check,topic,crit,statuses,testm,custom,kafka07,ft
,,15,2020-11-25T14:05:25.856575405Z,TESTING,rate-check,Rate Check,topic,crit,statuses,testm,custom,kafka07,ft
,,16,2020-11-25T14:05:25.576319539Z,PRODUCTION,rate-check,Rate Check,topic,ok,statuses,testm,custom,kafka01,ft
,,16,2020-11-25T14:05:25.856319359Z,TESTING,rate-check,Rate Check,topic,ok,statuses,testm,custom,kafka01,ft
,,16,2020-11-25T14:05:25.856444173Z,TESTING,rate-check,Rate Check,topic,ok,statuses,testm,custom,kafka01,ft
,,17,2020-11-25T14:05:25.856624989Z,TESTING,rate-check,Rate Check,topic,warn,statuses,testm,custom,kafka07,ft
"

topicsData = "
#group,false,false,false,false,true,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,double,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm,topic
,,0,2020-11-25T14:05:25.856485929Z,39.33716449678206,rate-check,Rate Check,KafkaMsgRate,crit,statuses,testm,custom,kafka07,ft,TESTING
,,0,2020-11-25T14:05:25.856575405Z,26.33716449678206,rate-check,Rate Check,KafkaMsgRate,crit,statuses,testm,custom,kafka07,ft,TESTING
,,1,2020-11-25T14:05:25.576319539Z,1.419231109050123,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka01,ft,PRODUCTION
,,2,2020-11-25T14:05:25.856319359Z,1.819231109049999,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka07,ft,TESTING
,,2,2020-11-25T14:05:25.856444173Z,1.635878190200181,rate-check,Rate Check,KafkaMsgRate,ok,statuses,testm,custom,kafka07,ft,TESTING
,,3,2020-11-25T14:05:25.856624989Z,8.33716449678206,rate-check,Rate Check,KafkaMsgRate,warn,statuses,testm,custom,kafka07,ft,TESTING

#group,false,false,false,false,true,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm,topic
,,4,2020-11-25T14:05:25.856485929Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 39.33716449678206,rate-check,Rate Check,_message,crit,statuses,testm,custom,kafka07,ft,TESTING
,,4,2020-11-25T14:05:25.856575405Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 26.33716449678206,rate-check,Rate Check,_message,crit,statuses,testm,custom,kafka07,ft,TESTING
,,5,2020-11-25T14:05:25.576319539Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert: ok - 1.419231109050123,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka01,ft,PRODUCTION
,,6,2020-11-25T14:05:25.856319359Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.819231109049999,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka07,ft,TESTING
,,6,2020-11-25T14:05:25.856444173Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.635878190200181,rate-check,Rate Check,_message,ok,statuses,testm,custom,kafka07,ft,TESTING
,,7,2020-11-25T14:05:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,rate-check,Rate Check,_message,warn,statuses,testm,custom,kafka07,ft,TESTING

#group,false,false,false,false,true,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,long,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm,topic
,,8,2020-11-25T14:05:25.856485929Z,1606313105623191313,rate-check,Rate Check,_source_timestamp,crit,statuses,testm,custom,kafka07,ft,TESTING
,,8,2020-11-25T14:05:25.856575405Z,1606313106696061106,rate-check,Rate Check,_source_timestamp,crit,statuses,testm,custom,kafka07,ft,TESTING
,,9,2020-11-25T14:05:25.576319539Z,1606313101487639160,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka01,ft,PRODUCTION
,,10,2020-11-25T14:05:25.856319359Z,1606313103477635916,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka07,ft,TESTING
,,10,2020-11-25T14:05:25.856444173Z,1606313104541635074,rate-check,Rate Check,_source_timestamp,ok,statuses,testm,custom,kafka07,ft,TESTING
,,11,2020-11-25T14:05:25.856624989Z,1606313107768317097,rate-check,Rate Check,_source_timestamp,warn,statuses,testm,custom,kafka07,ft,TESTING

#group,false,false,false,false,true,true,true,true,true,true,true,true,true,true
#datatype,string,long,dateTime:RFC3339,string,string,string,string,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,
,result,table,_time,_value,_check_id,_check_name,_field,_level,_measurement,_source_measurement,_type,host,realm,topic
,,12,2020-11-25T14:05:25.856485929Z,some detail: myrealm=ft,rate-check,Rate Check,details,crit,statuses,testm,custom,kafka07,ft,TESTING
,,12,2020-11-25T14:05:25.856575405Z,some detail: myrealm=ft,rate-check,Rate Check,details,crit,statuses,testm,custom,kafka07,ft,TESTING
,,13,2020-11-25T14:05:25.576319539Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka01,ft,PRODUCTION
,,14,2020-11-25T14:05:25.856319359Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka07,ft,TESTING
,,14,2020-11-25T14:05:25.856444173Z,some detail: myrealm=ft,rate-check,Rate Check,details,ok,statuses,testm,custom,kafka07,ft,TESTING
,,15,2020-11-25T14:05:25.856624989Z,some detail: myrealm=ft,rate-check,Rate Check,details,warn,statuses,testm,custom,kafka07,ft,TESTING
,,16,2020-11-25T14:05:25.856485929Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,crit,statuses,testm,custom,kafka07,ft,TESTING
,,16,2020-11-25T14:05:25.856575405Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,crit,statuses,testm,custom,kafka07,ft,TESTING
,,17,2020-11-25T14:05:25.576319539Z,Realm: ft - Hostname: kafka01 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka01,ft,PRODUCTION
,,18,2020-11-25T14:05:25.856319359Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka07,ft,TESTING
,,18,2020-11-25T14:05:25.856444173Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,ok,statuses,testm,custom,kafka07,ft,TESTING
,,19,2020-11-25T14:05:25.856624989Z,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,rate-check,Rate Check,id,warn,statuses,testm,custom,kafka07,ft,TESTING
"

// override source for existing topics
option tickscript._tssource = () =>
  csv.from(csv: topicsData)

dummy = array.from(rows: [{_time: 2020-01-01T00:00:00Z, _field: "uknown", _value: "unknown"}])

outData = "
#group,false,false,false,true,true,true,true,false,true,false,false,true,false,true,false,true,true
#datatype,string,long,double,string,string,string,string,string,string,long,dateTime:RFC3339,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,,,,
,result,table,KafkaMsgRate,_check_id,_check_name,_level,_measurement,_message,_source_measurement,_source_timestamp,_time,_type,details,host,id,realm,topic
,,0,1.819231109049999,rate-check,Rate Check,ok,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.819231109049999,testm,1606313103477635916,2020-11-25T14:05:25.856319359Z,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING
,,0,1.635878190200181,rate-check,Rate Check,ok,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.635878190200181,testm,1606313104541635074,2020-11-25T14:05:25.856444173Z,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING
,,1,39.33716449678206,rate-check,Rate Check,crit,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 39.33716449678206,testm,1606313105623191313,2020-11-25T14:05:25.856485929Z,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING
,,1,26.33716449678206,rate-check,Rate Check,crit,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 26.33716449678206,testm,1606313106696061106,2020-11-25T14:05:25.856575405Z,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING
,,2,8.33716449678206,rate-check,Rate Check,warn,statuses,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,testm,1606313107768317097,2020-11-25T14:05:25.856624989Z,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING
"

tickscript_from = (table=<-) =>
    tickscript.from(start: experimental.subDuration(d: 1m, from: now()), name: "TESTING")
      |> drop(columns: ["_start", "_stop"])

test _tickscript_from = () => ({
    input: dummy, // not really used but something is required for testing to work
	want: testing.loadMem(csv: outData),
	fn: tickscript_from,
})
