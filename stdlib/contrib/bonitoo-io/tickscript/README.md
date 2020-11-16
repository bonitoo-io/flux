# Tickscript package 

The `tickscript` package can be used to convert TICKscripts to InfluxDB tasks.

## tickscript.alert

`tickscript.alert()` checks input data and create alerts.

Alerts are records with state value other than OK or if the state just changed to OK status from a non OK state (ie. the alert recovered).

Parameters:
- `check` - _TODO_
- `id` - Function that constructs alert ID. Default is _TODO_
- `message` - Function that constructs alert message. Default is _TODO_
- `details` - Function that constructs detailed alert message.
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
- `notification` - _TODO_
- `endpoint` - Destination endpoint (eg. Slack endpoint).

## Examples
  
Alert task:

```js
import "contrib/bonitoo-io/tickscript"

// required task option
option task = {
  name: "Kafka Message Rate",
  every: 1m,
}

// Custom check info
check = {
  _check_id: "${task.name}-check",
  _check_name: "${task.name} Check",
  _type: "custom",
  tags: {},
} 

...

from(bucket: servicedb)
    |> range(start: -eval_duration)
    |> filterfn: (r) => ...)
    |> mean()
    |> duplicate(column: "_value", as: "KafkaMsgRate")
    |> group(columns: ["host", "realm"])
    |> tickscript.alert(
        check: check,
        id: (r) => "Realm: ${r.realm} - Hostname: ${r.host} / Metric: ${met_type} threshold alert",
        message: (r) => "${r.id}: ${r._level} - ${string(v:r.KafkaMsgRate)}",
        details: (r) => "https://grafana.nestlabs.com/dashboard/db/noc-jvm-type-realm-stats?&panelId=19&fullscreen&orgId=1&var-myrealm=${r.realm}",
        crit: (r) => r.KafkaMsgRate > h_threshold or r.KafkaMsgRate < l_threshold,
    )
    |> tickscript.topic(name: sltest)
```

Topic handler task:

```js
import "contrib/bonitoo-io/tickscript"
import "slack"

// required task option
option task = {
  name: "Noc Testing Topic",
  every: 1m,
  topic: "NOC_TESTING",
}

// Custom notification rule
notification = {
  _notification_rule_id: "${task.topic}-rule",
  _notification_rule_name: "${task.name} Rule",
  _notification_endpoint_id: "${task.topic}-endpoint",
  _notification_endpoint_name: "${task.name} Endpoint",
}

slack_endpoint = slack.endpoint(url: "https://hooks.slack.com/services/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")(mapFn: (r) => ({
    channel: "",
    text: "Message: ${r._message}\n\nDetail: ${r.details}",
    color: if r._level == "ok" then "good" else "warning"
}))

tickscript.from(start: -task.every, name: task.topic)
    |> tickscript.notify(notification: notification, endpoint: slack_endpoint)
```