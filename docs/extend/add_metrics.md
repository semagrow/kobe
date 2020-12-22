# Add new metrics

In this section we describe the process of adding new metrics in KOBE.

## Step 1. Produce a log message that contains the metric

To add a new metric you must first modify your compoment to produce it.
KOBE collects metrics by parsing log messages so to add the metric just add another
log message in the compoment for this metric.

Except from the metric, your log message should contain the following information:

* Experiment name
* Start time of the experiment
* Query name
* Run

If the log message is produced by an evaluator, these parameters are obtained
from the KOBE operator (see [here](./add_evaluator.md#step-1-prepare-your-docker-image)).
If the log message is produced by a federator, these parameters are obtained
by a SPARQL comment (see [here](./support_metrics.md#step-1-provide-support-for-all-evaluation-metrics)).

## Step 2. Configure Fluentd with a corresponding regex pattern

All metrics are captured by Fluentd, so to collect the new metric you must edit the
`operator/deploy/efk-config/fluentd-values.yaml` file and add
another Fluentd filter at `containers.input.conf` for the introduced
log pattern of the metric.

For example if the log pattern is

`^Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - MyMetric: (?<my_metric>[0-9]+)$`

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
    expression ^Experiment: (?<experiment>[^ ]+) - Date: (?<date>[^ ]+ [^ ]+) - Query: (?<query>[^ ]+) - Run: (?<run>[0-9]+) - MyMetric: (?<my_metric>[0-9]+)$
    types my_metric:integer
  </parse>
</filter>
```

After redeployling Fluentd, every time the metric log message is produced, it will be captured by Fluentd and stored in ElasticSearch.


## Step 3. Configure Kibana to visualize your metric

Finally, you must add a [new Kibana Dashboard](https://www.elastic.co/guide/en/kibana/7.x/dashboard.html) to visualize the new metric, or add a Visualization on one of the [existing Dashboards](https://www.elastic.co/guide/en/kibana/7.x/edit-dashboards.html).

## Example

As an example, check out the
[Fluentd configuration](https://github.com/semagrow/kobe/tree/devel/operator/deploy/efk-config/fluentd-values.yaml)
for the metrics we currently support.
