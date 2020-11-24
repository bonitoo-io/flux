# Tickscript package

The `tickscript` package can be used to convert TICKscripts to InfluxDB tasks.

## Available functions

- `alert`
- `topic`
- `from`
- `notify`
- `as`

Many TICKscript functions has similar counterparts in Flux.

## Conversion rules

* Both `batch` and `stream` in TICKscript translates to `from(bucket: ...)` in Flux.

* `every(duration)` property maps to task's option `every` field.
  For better control or aligned scheduling, use `cron` option instead.

* `period(duration)` property maps to `range(start: -duration)` in Flux pipeline.

* `groupBy(columns)` maps to `group(columns)`.
  Columns must include internal `_measurement` column.

* TICKscript `alert` provides property methods to send alerts to event handlers or a topic.
  In Flux, use `tickscript.notify()` or `tickscript.topic()` pipeline functions.

* TICKscript pipeline with multiple alerts translates to multiple Flux pipelines, ie.

```js
var data = batch
    | query(...)
data
    | alert()
        .topic('A')
    | alert()
        .topic('B')
```

becomes

```js
data = from(bucket: ...)
    |> range(start: -duration)
    ...
data
    |> alert()
    |> topic('A')
data
    |> alert()
    |> topic('B')
```

## tickscript.alert

`tickscript.alert()` checks input data and create alerts.
It requires pivoted data (call `schema.fieldsAsCols()` before `tickscript.alert()`).

Alerts are records with state value other than OK or if the state just changed to OK status from a non OK state (ie. the alert recovered).

Parameters:
- `check` - Required by underlying `monitor` package _TODO_
- `id` - Function that constructs alert ID. Default is _TODO_
- `message` - Function that constructs alert message. Default is _TODO_
- `details` - Function that constructs detailed alert message. Default is _TODO_
- `crit` - Predicate function that determines `crit` status. Default is `(r) => false`.
- `warn` - Predicate function that determines `warn` status. Default is `(r) => false`.
- `info` - Predicate function that determines `info` status. Default is `(r) => false`.
- `ok` - Predicate function that determines `info` status. Default is `(r) => true`.
- `stateChangesOnly` - Only send alerts where the state changed. Default value: `false`.

## tickscript.topic

`tickscript.topic()` sends alerts to a topic.

Parameters:
- `name` - Topic name.

## tickscript.from

`tickscript.from()` reads alerts from a topic.

Parameters:
- `name` - Topic name.
- `start` - Time range start.
- `stop` - Time range stop. Default value is `now()`.
- `fn` - Predicate function. Default value is `fn=(r) => true`.

## tickscript.notify

`tickscript.notify()` sends alerts to an endpoint.

Parameters:
- `notification` - Required by underlying `monitor` package _TODO_
- `endpoint` - Destination endpoint (eg. Slack endpoint).

## tickscript.as

`tickscript.as()` is a convenience function for renaming a column.

It is equivalent to using `rename(fn: (column) => ...)`

Parameters:
- `column` - Existing column. Default value is `_value`.
- `as` - Desired column name.

## Examples

### Using topic

TICKscript:
```js
metric_type = 'kafka_message_in_rate'
...

batch
    |query('SELECT ' + metric_type + ' AS "KafkaMsgRate" FROM ' +  db  + ' WHERE realm = \'' + tier + '\' AND "host" =~ /^kafka.+.m02/')
    .period(duration)
    .every(frequency)
    .groupBy('host','realm')
   |alert()
        .id('Realm: {{index .Tags "realm"}} - Hostname: {{index .Tags "host"}} / Metric: ' + metric_type + ' threshold alert' )
        .message('{{.ID }}: {{ .Level }} - {{ index .Fields "KafkaMsgRate" | printf "%0.2f"}}')
        .crit(lambda: ("KafkaMsgRate" > h_threshold) OR ("KafkaMsgRate" < l_threshold) )
        .stateChangesOnly()
        .topic('TESTING')
```

InfluxDB alert task:

```js
import "contrib/bonitoo-io/tickscript"
import "influxdata/influxdb/schema"

// required task option
option task = {
  name: "Kafka Message Rate",
  every: 1m,
}

// custom check info
check = {
  _check_id: "${task.name}-check",
  _check_name: "${task.name} Check",
  _type: "custom",
  tags: {},
} 

// variables
metric_type = 'kafka_message_in_rate'
...

from(bucket: servicedb)
    |> range(start: -period)
    |> filter(fn: (r) => r._field == metric_type and r.realm == tier and r.host =~ /^kafka.+.m02/)
    |> group(columns: ["_measurement", "host", "realm"])
    |> schema.fieldsAsCols()
    |> drop(columns: ["_start", "_stop"])
    |> tickscript.as(column: metric_type, as: "KafkaMsgRate")
    |> tickscript.alert(
        check: check,
        id: (r) => "Realm: ${r.realm} - Hostname: ${r.host} / Metric: ${metric_type} threshold alert",
        message: (r) => "${r.id}: ${r._level} - ${string(v:r.KafkaMsgRate)}",
        crit: (r) => r.KafkaMsgRate > h_threshold or r.KafkaMsgRate < l_threshold,
        stateChangeOnly: true
    )
    |> tickscript.topic(name: "TESTING")
```

Topic handler task:

```js
import "contrib/bonitoo-io/tickscript"
import "slack"

// required task option
option task = {
  name: "Testing Topic Handler",
  every: 1m,
}

// custom notification rule
notification = {
  _notification_rule_id: "${task.topic}-rule",
  _notification_rule_name: "${task.name} Rule",
  _notification_endpoint_id: "${task.topic}-endpoint",
  _notification_endpoint_name: "${task.name} Endpoint",
}

// destination endpoint
slack_endpoint = slack.endpoint(url: "https://hooks.slack.com/services/...")(mapFn: (r) => ({
    channel: "",
    text: "Message: ${r._message}\n\nDetail: ${r.details}",
    color: if r._level == "ok" then "good" else "warning"
}))

tickscript.from(start: -task.every, name: "TESTING")
    |> tickscript.notify(notification: notification, endpoint: slack_endpoint)
```

### Sending alerts directly to event handler

Task:

```js
import "contrib/bonitoo-io/tickscript"
import "influxdata/influxdb/schema"
import "slack"

// required task option
option task = {
  name: "Kafka Message Rate",
  every: 1m,
}

// custom check info
check = {
  _check_id: "${task.name}-check",
  _check_name: "${task.name} Check",
  _type: "custom",
  tags: {},
}

// custom notification rule
notification = {
  _notification_rule_id: "${task.topic}-rule",
  _notification_rule_name: "${task.name} Rule",
  _notification_endpoint_id: "${task.topic}-endpoint",
  _notification_endpoint_name: "${task.name} Endpoint",
}

// variables
metric_type = 'kafka_message_in_rate'
...

// destination endpoint
slack_endpoint = slack.endpoint(url: "https://hooks.slack.com/services/...")(mapFn: (r) => ({
    channel: "",
    text: "Message: ${r._message}\n\nDetail: ${r.details}",
    color: if r._level == "ok" then "good" else "warning"
}))

from(bucket: servicedb)
    |> range(start: -period)
    |> filter(fn: (r) => r._field == met_type and r.realm == tier and r.host =~ /^kafka.+.m02/)
    |> group(columns: ["_measurement", "host", "realm"])
    |> schema.fieldsAsCols()
    |> drop(columns: ["_start", "_stop"])
    |> tickscript.as(column: metric_type, as: "KafkaMsgRate")
    |> tickscript.alert(
        check: check,
        id: (r) => "Realm: ${r.realm} - Hostname: ${r.host} / Metric: ${metric_type} threshold alert",
        message: (r) => "${r.id}: ${r._level} - ${string(v:r.KafkaMsgRate)}",
        crit: (r) => r.KafkaMsgRate > h_threshold or r.KafkaMsgRate < l_threshold,
        stateChangeOnly: true
    )
    |> tickscript.notify(notification: notification, endpoint: slack_endpoint)
```

## TODO

* provide helper functions for instantiating `check` and `notification` custom records
