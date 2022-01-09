# Owl

- [Owl](#owl)
  - [Job](#job)
    - [HTTP Job](#http-job)
    - [File Job](#file-job)
  - [Flag](#flag)
  - [Config](#config)
  - [Config Example](#config-example)
  - [Schema](#schema)
    - [Job Header Document](#job-header-document)
    - [CRON Config Object](#cron-config-object)
    - [HTTP Job Config Document](#http-job-config-document)
    - [File Job Config Document](#file-job-config-document)
    - [Job Parameter Object](#job-parameter-object)
    - [Job Parameter Value List](#job-parameter-value-list)

Owl is a server to schedule jobs.
Owl takes advantage of CRON expression to config jobs.

Owl supports:

- [HTTP Job](#http-job)
- [File Job](#file-job)

## Job

### HTTP Job

Each HTTP job will send a POST request to the specified URI,
headers of the request are:

| Header | Value |
| --- | --- |
| Content-Type | application/json |

and the config of [Job Parameter Object](#job-parameter-object) will be parsed as a JSON-form request body.

For example, there is a config like:

```yaml
---
name: example-http-job
type: http
cron:
  express: "* * * * * *"
  skip_if_still_running: false
  delay_if_still_running: true
---
uri: http://localhost/example
parameters:
  - name: rid
    value: $RANDOM_ID
  - name: name
    value: tester
  - name: param
    value: "123"
```

parameters will be parse as JSON like:

```json
{
  "rid": "$RANDOM_ID",
  "name": "tester",
  "param": "123"
}
```

and send to `http://localhost/example`.

**CAUTION**: If there are 2 or above of the **SAME** name of [Job Parameter Object](#job-parameter-object), Owl would take the **LAST** one value.

### File Job

File job runs an executable file with flags(if any).

The config of [Job Parameter Object](#job-parameter-object) will be parsed as flags.

For example, there is a config like:

```yaml
---
name: example-files-job
type: files
cron:
  express: "* * * * * *"
  skip_if_still_running: false
  delay_if_still_running: true
---
path: ./example.exe
parameters:
  - name: --rid
    value: $RANDOM_ID
  - name: --name
    value: tester
  - name: --development
    value: ""
```

and the config will be parsed like running on the command line as below:

```cmd
./example.exe --rid=c79e9307lsqjbn6vbmd0 --name=tester --development
```

## Flag

All supported flags are as below:

| Flag | Description |
| --- | --- |
| --dev | |
| --dir | The directory to load configurations. The way to create configurations, see [Config](#config). |

## Config

Every configurations are [YAML](https://en.wikipedia.org/wiki/YAML) form.
There are two documents in the configuration file:

| Index Of Document | Type | Description |
| --- | --- | --- |
| 0 | [Job Header Document](#job-header-document) | **REQUIRED**. |
| 1 | [HTTP Job Config Document](#http-job-config-document) \| [File Job Config Document](#file-job-config-document) | **REQUIRED**. Define how to execute the job. The type of the document depends on the value of the `type` field in [index-0 document](#job-header-document). |

## Config Example

```yaml
---
name: example-http-job
type: http
cron:
  express: "* * * * *"
  skip_if_still_running: false
  delay_if_still_running: true
---
uri: http://localhost/example
parameters:
  - name: rid
    value: $RANDOM_ID
```

## Schema

### Job Header Document

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| name | string | **REQUIRED**. The name of the job. |
| type | string | **REQUIRED**. The available values are '*http*', '*files*'. |
| cron | [CRON Config Object](#cron-config-object) | The CRON setting. |

### CRON Config Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| express | string | **REQUIRED**. An expression represents a set of times, using 6 space-separated fields like '`* * * * * *`'. See [CRON Expression Format](https://pkg.go.dev/github.com/robfig/cron#hdr-CRON_Expression_Format) for details. |
| skip_if_still_running | boolean | If true, not to start the job util the running job is done. |
| delay_if_still_running | boolean | If true, wait util the previous job is done and start. |

### HTTP Job Config Document

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| uri | string | **REQUIRED**. |
| parameters | \[\][Job Parameter Object](#job-parameter-object) | Parameters are sent as a JSON-form request body. |

### File Job Config Document

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| path | string | **REQUIRED**. An executable file path. |
| parameters | \[\][Job Parameter Object](#job-parameter-object) | Set the given flags if any. |

### Job Parameter Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| name | string | **REQUIRED**. The given name of the parameter. |
| value | string | **REQUIRED**. See [Job Parameter Value List](#job-parameter-value-list) for details. |

### Job Parameter Value List

The value with prefix '*$*' is reserved for the system value as the following table:

| value | comment |
| --- | --- |
| `$RANDOM_ID` | A 20-length of the random value. It is useful to track the session scope. |
