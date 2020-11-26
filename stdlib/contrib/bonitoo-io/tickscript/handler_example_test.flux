package tickscript_test

import "testing"
import "csv"
import ts "contrib/bonitoo-io/tickscript"
import "experimental"
import "influxdata/influxdb/monitor"

option now = () => (2020-11-25T14:06:00Z)

option monitor.log = (tables=<-) => tables

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
option ts._tssource = () =>
  csv.from(csv: topicsData)

outData = "
#group,false,false,false,true,true,false,true,false,true,true,true,true,true,false,false,true,false,true,false,true,true,true
#datatype,string,long,double,string,string,string,string,string,string,string,string,string,string,long,long,string,string,string,string,string,string,string
#default,_result,,,,,,,,,,,,,,,,,,,,,
,result,table,KafkaMsgRate,_check_id,_check_name,_level,_measurement,_message,_notification_endpoint_id,_notification_endpoint_name,_notification_rule_id,_notification_rule_name,_source_measurement,_source_timestamp,_status_timestamp,_type,details,host,id,realm,topic,_sent
,,0,1.819231109049999,rate-check,Rate Check,ok,notifications,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.819231109049999,rate-endpoint,Rate Endpoint,rate-rule,Rate Rule,testm,1606313103477635916,1606313125856319359,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING,true
,,0,1.635878190200181,rate-check,Rate Check,ok,notifications,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: ok - 1.635878190200181,rate-endpoint,Rate Endpoint,rate-rule,Rate Rule,testm,1606313104541635074,1606313125856444173,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING,true
,,0,39.33716449678206,rate-check,Rate Check,crit,notifications,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 39.33716449678206,rate-endpoint,Rate Endpoint,rate-rule,Rate Rule,testm,1606313105623191313,1606313125856485929,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING,true
,,0,26.33716449678206,rate-check,Rate Check,crit,notifications,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: crit - 26.33716449678206,rate-endpoint,Rate Endpoint,rate-rule,Rate Rule,testm,1606313106696061106,1606313125856575405,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING,true
,,0,8.33716449678206,rate-check,Rate Check,warn,notifications,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert: warn - 8.33716449678206,rate-endpoint,Rate Endpoint,rate-rule,Rate Rule,testm,1606313107768317097,1606313125856624989,custom,some detail: myrealm=ft,kafka07,Realm: ft - Hostname: kafka07 / Metric: kafka_message_in_rate threshold alert,ft,TESTING,true
"

notification = {
  _notification_rule_id: "rate-rule",
  _notification_rule_name: "Rate Rule",
  _notification_endpoint_id: "rate-endpoint",
  _notification_endpoint_name: "Rate Endpoint",
}

endpoint = () => (tables=<-) => tables |> experimental.set(o: { _sent: "true" })

tickscript_handler_example = (table=<-) =>
  ts.from(start: experimental.subDuration(d: 1m, from: now()), name: "TESTING")
    |> ts.notify(notification: notification, endpoint: endpoint())
    |> drop(columns: ["_start", "_stop", "_time"])

test _tickscript_handler_example = () => ({
    input: testing.loadStorage(csv: topicsData), // not really used but something is required for testing to work
	want: testing.loadMem(csv: outData),
	fn: tickscript_handler_example,
})
