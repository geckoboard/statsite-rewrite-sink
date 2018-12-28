# Statsite rewrite sink

A [statsite sink](https://github.com/statsite/statsite) that rewrites
metric names before passing the metrics through to another statsite
sink.

This was born from a need to extract tags from consul/vault/envoy
metrics so that could make better use of librato's support for tagged
metrics.

This is very similar to [seatgeek's
statsd-rewrite-proxy](https://github.com/seatgeek/statsd-rewrite-proxy),
the main difference being that their proxy works on the raw statsd
packets sent from applications, and this project works on aggregated
metrics statsite flushes to sinks/backends/integrations. This takes the
rewrite logic out of the hot-path (statsite flushes its metrics in a
separate process from the main thread) and means we can avoid messing
around with UDP packets.

The other difference is that we eventually want to support extracting
tags from part of a metric name, (envoy [has example rules for
extracting tags from parts of a metric
name](https://github.com/envoyproxy/envoy/blob/87553968ec2258919d986e9a76512b0009d01575/source/common/config/well_known_names.cc))
