# Add new metrics

In this section we describe the process of adding new metrics in KOBE.

To add a new metric you must first modify your compoment to produce it. KOBE collects metrics by parsing log messages 
so to add the metric just add another log message in the compoment for this metric.

All metrics are captured by Fluentd, so to collect the new metric you must edit the `operator/deploy/efk-config/fluentd-values.yaml` file and add
another Fluentd filter at `containers.input.conf` for the introduced log pattern of the metric.

For example if the log pattern is

`^I - [^ ]+ [^ ]+ - .{12} - .{20} - [^ ]+ - Experiment: (?<experiment_date>[^ ]+ - Date: [^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - Query Evaluation Time: (?<evaluation_time>[0-9]+) - Results: (?<results>[0-9]+)$`

you must add filter

```
<filter kubernetes.**>
  @type parser
  key_name message
  reserve_time true
  reserve_data true
  #suppress_parse_error_log true
  <parse>
    @type regexp
    expression ^I - [^ ]+ [^ ]+ - .{12} - .{20} - [^ ]+ - Experiment: (?<experiment_date>[^ ]+ - Date: [^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - Query Evaluation Time: (?<evaluation_time>[0-9]+) - Results: (?<results>[0-9]+)$
    types evaluation_time:integer,results:integer
  </parse>
</filter>
```

After redeployling Fluentd, every time the metric log message is produced, it will be captured by Fluentd and stored in ElasticSearch.

Finally, you must add a [new Kibana Dashboard](https://www.elastic.co/guide/en/kibana/7.x/dashboard.html) to visualize the new metric, or add a Visualization on one of the [existing Dashboards](https://www.elastic.co/guide/en/kibana/7.x/edit-dashboards.html).
