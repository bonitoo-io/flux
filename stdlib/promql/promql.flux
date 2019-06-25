package promql

builtin changes
builtin dayOfMonth
builtin dayOfWeek
builtin daysInMonth
builtin emptyTable
builtin extrapolatedRate
builtin hour
builtin instantRate
builtin minute
builtin month
builtin resets
builtin timestamp
builtin year

// hack to simulate an imported promql package
promql = {
  dayOfMonth:dayOfMonth,
  dayOfWeek:dayOfWeek,
  daysInMonth:daysInMonth,
  hour:hour,
  minute:minute,
  month:month,
  year:year,
}
