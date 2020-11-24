package tickscript

import "experimental"
import "influxdata/influxdb"
import "influxdata/influxdb/monitor"
import "influxdata/influxdb/schema"

option bucket = "kapacitor"
option rp = 7d

// override monitor persistence functions to use our bucket instead of "_monitoring"
option monitor.write = (tables=<-) =>
  tables |> experimental.to(bucket: bucket)
option monitor.log = (tables=<-) =>
  tables |> experimental.to(bucket: bucket)

// removes column from group key
_ungroup = (column, tables=<-) =>
  tables
    |> duplicate(column: column, as: "____temp_column____")
    |> drop(columns: [column])
    |> rename(columns: {____temp_column____: "_level"})

// sorts by columns with handy defaults
_sort = (columns=["_source_timestamp"], desc=false, tables=<-) =>
  tables
    |> sort(columns: columns, desc: desc)

// last statuses (one per series)
_last = (start, stop=now()) =>
  influxdb.from(bucket: bucket)
    |> range(start: start, stop: stop)
    |> schema.fieldsAsCols()
    |> drop(columns: ["_start", "_stop"])
    |> _ungroup(column: "_level")
    |> _sort()
    |> last(column: "_source_timestamp")
    |> experimental.group(mode: "extend", columns: ["_level"])

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
  lastStatuses = _last(start: -rp)
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
  sequence = union(tables: [lastStatuses, statuses])
  notOk = statuses
    |> filter(fn: (r) => r._level != "ok")
  againOk = sequence
    |> monitor.stateChanges(toLevel: "ok")
  anyChange = sequence
    |> monitor.stateChangesOnly()
  alerts =
    if stateChangesOnly then
      anyChange
    else
      union(tables: [notOk, againOk])
  return alerts
}

// routes alerts to topic
topic = (name, tables=<-) =>
  tables
    |> map(fn: (r) => ({ r with topic: name }))
    // use this to have extra series (topic is a tag), otherwise it overrides existing statuses (topic is a field)
    // |> experimental.group(mode: "extend", columns: ["topic"])
    |> experimental.to(bucket: bucket)

// sends alerts to event handler
notify = (notification, endpoint, tables=<-) =>
  tables
    |> monitor.notify(data: notification, endpoint: endpoint)

// reads topic
from = (name, start, stop=now(), fn=(r) => true) =>
  influxdb.from(bucket: bucket)
    |> range(start: start, stop: stop)
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
    |> _ungroup(column: "_level")
    |> _sort()

// renames column
// it is meant to be a convenience function to rename result column when "SELECT x AS y" is used in TICKscript and x is variable
as = (column="_value", as, tables=<-) => {
  _column = column
  _as = as
  return
    tables
      |> rename(fn: (column) => if column == _column then _as else column)
}

// groups by specified columns
// it is meant to be a convenience function, it adds _measurement column which is required by monitor.check()
groupBy = (columns, tables=<-) =>
  tables
    |> group(columns: columns)
    |> experimental.group(columns: ["_measurement"], mode:"extend") // required by monitor.check
