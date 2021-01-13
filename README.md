[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/check-load)
![Go Test](https://github.com/nixwiz/check-load/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/check-load/workflows/goreleaser/badge.svg)

# Sensu Load Average Check

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
  - [Help output](#help-output)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Check definition](#check-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu Load Average Check is a [Sensu Check][1] that provides alerting and
metrics for [load average][2].  Metrics are provided in [nagios perfdata][3]
format.

This check is available for Linux and macOS only.

## Usage examples

The critical and warning thresholds are based on multipliers of the number of
CPUs present on a system.  For example, when using all defaults the critical
threshold is 2 times the number of physical CPUs. This means for a system with
4 CPUs a critical event would occur when the 1 minute load average is equal to
or greater than 8. In this same scenario, a warning event would occur if the 1
minute load average is less than 8 but greater than or equal to 6, given that
the warning multiplier is 1.5.

### Help output

```
Sensu Load Average Check

Usage:
  check-load [flags]
  check-load [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -c, --critical-multiplier float   The critical threshold multiplier (# CPUs x multiplier) (default 2)
  -w, --warning-multiplier float    The warning threshold multiplier (# CPUs x multiplier) (default 1.5)
  -a, --compare-all-intervals       Compare thresholds to all (1m, 5m, 15m) load averages
  -l, --count-logical-cpu           Include Logical CPUs (e.g. hyperthreading) in factoring thresholds
  -h, --help                        help for check-load

Use "check-load [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][4] are the best way to make use of this plugin. If you're not
using an asset, please consider doing so! If you're using sensuctl 5.13 with
Sensu Backend 5.13 or later, you can use the following command to add the asset:

```
sensuctl asset add nixwiz/check-load
```

If you're using an earlier version of sensuctl, you can find the asset on the
[Bonsai Asset Index][5].

### Check definition

```yml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-load
  namespace: default
spec:
  command: >-
    check-load
    --warning-multiplier 2
    --critical-multiplier 3
  subscriptions:
  - system
  runtime_assets:
  - nixwiz/check-load
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an
Asset. If you would like to compile and install the plugin from source or
contribute to it, download the latest version or create an executable from this
source.

From the local path of the check-load repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][6].

[1]: https://docs.sensu.io/sensu-go/latest/reference/checks/
[2]: http://www.brendangregg.com/blog/2017-08-08/linux-load-averages.html
[3]: https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-schedule/collect-metrics-with-checks/#supported-output-metric-formats
[4]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[5]: https://bonsai.sensu.io/assets/nixwiz/check-load
[6]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
