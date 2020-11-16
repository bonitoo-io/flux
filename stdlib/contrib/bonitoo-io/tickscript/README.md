# Tickscript package 

The `tickscript` package can be used to convert TICKscripts to InfluxDB tasks.

## tickscript.alert

_TODO_

## tickscript.topic

_TODO_

## tickscript.from

_TODO_

## tickscript.notify

_TODO_

## Examples
  
Alert task:

```javascript
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
        id: (r) => "Realm: ${r.realm} - Hostname: ${r.host} / Metric: ${met_type} threshold alert",
        message: (r) => "${r.id}: ${r._level} - ${string(v:r.KafkaMsgRate)}",
        details: (r) => "https://grafana.nestlabs.com/dashboard/db/noc-jvm-type-realm-stats?&panelId=19&fullscreen&orgId=1&var-myrealm=${r.realm}",
        crit: (r) => r.KafkaMsgRate > h_threshold or r.KafkaMsgRate < l_threshold,
    )
    |> tickscript.topic(name: pgtest)
```

Topic handler task:

```javascript
import "contrib/bonitoo-io/tickscript"
import "slack"

// required task option
option task = {
  name: "Noc Testing Topic",
  topic: "NOC_TESTING",
  every: 1m,
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