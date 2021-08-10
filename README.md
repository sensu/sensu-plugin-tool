# Sensu Plugin Tool

- [Overview](#overview)
- [Usage examples](#usage-examples)
  - [Check](#check-example)
  - [Handler](#handler-example)
  - [AWS Service](#aws-service-example)
- [Installing from source and
  contributing](#installation-from-source-and-contributing)

## Overview

The [Sensu Plugin Tool][0] is a command-line tool to generates scaffolding for a
new [Sensu][1] Plugin project.

The following plugin types are currently supported:

* Check - [default check template][2]
* Handler - [default handler template][3]
* Mutator - [default mutator template][4]
* Sensuctl - [default sensuctl template][5]
* Custom - Must specify plugin template

### Custom plugin templates
While the default templates are sufficient to get you started, there are specialized
Cloud service templates available that maybe a better starting point when building
new plugins targeting a cloud provider.

Currently available service template include:

* AWS - [aws service check template][7]

## Usage examples

### Check example
Creating a check plugin using interactive mode:

```sh
$ sensu-plugin-tool new check
? Template URL https://github.com/sensu/check-plugin-template
? Project name My Check
? Description Description for My Check
? Github User githubuser
? Github Project my-check
? Copyright Year 2020
? Copyright Holder Me
Success!
```

Creating a check plugin using flags:

```sh
sensu-plugin-tool new check \
    --template-url "https://github.com/sensu/check-plugin-template" \
    --name "My Check" \
    --description "Description for My Check" \
    --github-user mygithubuser \
    --github-project my-check \
    --copyright-year 2020 \
    --copyright-holder Me
```


### Handler example
Creating a handler plugin using interactive mode:

```sh
$ sensu-plugin-tool new handler
? Template URL https://github.com/sensu/handler-plugin-template
? Project name My Handler
? Description Description for My Handler
? Github User githubuser
? Github Project my-handler
? Copyright Year 2020
? Copyright Holder Me
Success!
```

Creating a handler plugin using flags:

```sh
sensu-plugin-tool new handler \
    --template-url "https://github.com/sensu/handler-plugin-template" \
    --name "My Handler" \
    --description "Description for My Handler" \
    --github-user mygithubuser \
    --github-project my-handler \
    --copyright-year 2020 \
    --copyright-holder Me
```

### AWS service example
Creating an AWS service plugin using interactive mode:

```sh
$ sensu-plugin-tool new custom
? Template URL https://github.com/sensu/aws-plugin-template
? Project name My AWS Service Plugin
? Description Description for My AWS Service Plugin
? Github User githubuser
? Github Project my-aws-service-plugin
? Copyright Year 2020
? Copyright Holder Me
Success!
```

Creating an AWS service plugin using flags:

```sh
sensu-plugin-tool new custom \
    --template-url "https://github.com/sensu/aws-plugin-template" \
    --name "My AWS Service Plugin" \
    --description "Description for My AWS Service Plugin" \
    --github-user mygithubuser \
    --github-project my-aws-service-plugin \
    --copyright-year 2020 \
    --copyright-holder Me
```

## Installing from source and contributing

Download the latest version of the sensu-plugin-tool from [releases][6],
or create an executable script from this source.

### Compiling

From the local path of the sensu-plugin-tool repository:

``` sh
go build
```

[0]: https://github.com/sensu/sensu-plugin-tool
[1]: https://sensu.io
[2]: https://github.com/sensu/check-plugin-template
[3]: https://github.com/sensu/handler-plugin-template
[4]: https://github.com/sensu/mutator-plugin-template
[5]: https://github.com/sensu/sensuctl-plugin-template
[6]: https://github.com/sensu/sensu-plugin-tool/releases
[7]: https://github.com/sensu/aws-plugin-template
