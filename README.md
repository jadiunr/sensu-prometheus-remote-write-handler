[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/jadiunr/sensu-prometheus-remote-write-handler)
![goreleaser](https://github.com/jadiunr/sensu-prometheus-remote-write-handler/workflows/goreleaser/badge.svg)
[![Go Test](https://github.com/jadiunr/sensu-prometheus-remote-write-handler/workflows/Go%20Test/badge.svg)](https://github.com/jadiunr/sensu-prometheus-remote-write-handler/actions?query=workflow%3A%22Go+Test%22)
[![goreleaser](https://github.com/jadiunr/sensu-prometheus-remote-write-handler/workflows/goreleaser/badge.svg)](https://github.com/jadiunr/sensu-prometheus-remote-write-handler/actions?query=workflow%3Agoreleaser)

# Sensu Prometheus remote write Handler

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
  - [Help output](#help-output)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
- [Installation from source](#installation-from-source)
- [Contributing](#contributing)

## Overview

The Sensu Prometheus remote write Handler is a [Sensu Handler][6] that sends metrics to time series database that has Prometheus remote write interface.

## Usage examples

### Help output

```
Prometheus remote write Handler for Sensu

Usage:
  sensu-prometheus-remote-write-handler [flags]
  sensu-prometheus-remote-write-handler [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -e, --endpoint string   Remote write endpoint
  -H, --header strings    Additional header(s) to send in remote write request
  -h, --help              help for sensu-prometheus-remote-write-handler
  -t, --timeout string    Remote write timeout (default "10s")

Use "sensu-prometheus-remote-write-handler [command] --help" for more information about a command.
```

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add jadiunr/sensu-prometheus-remote-write-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/jadiunr/sensu-prometheus-remote-write-handler].

### Handler definition

```yml
---
type: Handler
api_version: core/v2
metadata:
  name: sensu-prometheus-remote-write-handler
  namespace: default
spec:
  command: sensu-prometheus-remote-write-handler -e http://localhost:9009/prometheus -H "X-Scope-OrgID:tenant-example"
  type: pipe
  runtime_assets:
  - jadiunr/sensu-prometheus-remote-write-handler
```

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the sensu-prometheus-remote-write-handler repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu/handler-plugin-template/blob/master/.github/workflows/release.yml
[5]: https://github.com/sensu/handler-plugin-template/actions
[6]: https://docs.sensu.io/sensu-go/latest/reference/handlers/
[7]: https://github.com/sensu/handler-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
