package tickscript

import "experimental"
import "influxdata/influxdb"
import "influxdata/influxdb/monitor"
import "influxdata/influxdb/schema"

option bucket = "kapacitor"

// override monitor persistence functions to use our bucket instead of "_monitoring"
option monitor.write = (tables=<-) =>
  tables |> experimental.to(bucket: bucket)
option monitor.log = (tables=<-) =>
  tables |> experimental.to(bucket: bucket)

// filters statuses that represent Kapacitor alerts, ie. non-OK statuses and non-OK to OK state changes
_alertsFromStatuses = (tables=<-) => {
    notOk = tables
        |> filter(fn: (r) => r._level != "ok")
    againOk = tables
        |> monitor.stateChanges(toLevel: "ok")
    return
        union(tables: [notOk, againOk])
}

// alert
alert = (
    check,
    id=(r)=>"${r._measurement}:TODO",
    message=(r)=>"${r.id} is ${r._level}",
    details=(r)=>"TODO",
    crit=(r) => false,
    warn=(r) => false,
    info=(r) => false,
    ok=(r) => true,
    stateChangesOnly=false,
    tables=<-) => {
  statuses = tables
    |> map(fn: (r) => ({ r with id: id(r: r) }))
    |> map(fn: (r) => ({ r with details: details(r: r) }))
    |> monitor.check(
        crit: crit,
        warn: warn,
        info: info,
        ok: ok,
        messageFn: message,
        data: check
    )
  return
    if stateChangesOnly then
      statuses
        |> monitor.stateChangesOnly()
    else
      statuses
        |> _alertsFromStatuses()
}

// routes alerts to topic
topic = (name, tables=<-) =>
  tables
    |> map(fn: (r) => ({ r with topic: name }))
    // use this to have extra series, otherwise it overrides existing statuses
    // |> experimental.group(mode: "extend", columns: ["topic"])
    |> experimental.to(bucket: bucket)

// sends alerts to event handler
notify = (notification, endpoint, tables=<-) =>
  tables
    |> monitor.notify(data: notification, endpoint: endpoint)

// reads topic
from = (name, start, stop=now(), fn=(r) => true) =>
  influxdb.from(bucket: bucket)
    |> range(start: start)
    |> filter(fn: (r) => r._measurement == "statuses")
    // ### use this when topic is in the group key (see topic function)
    // |> filter(fn: (r) => r.topic == topic)
    // |> filter(fn: fn)
    // ###
    |> schema.fieldsAsCols()
    // ### use this when topic is not in the group key (see topic function)
    |> filter(fn: (r) => r.topic == name)
    |> filter(fn: fn)
    // ###
