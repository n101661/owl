# Owl

## Config

Every configurations are [YAML](https://en.wikipedia.org/wiki/YAML) form.
There are two documents in the configuration file:

| Index Of Document | Type | Description |
| --- | --- | --- |
| 0 | [Job Header Document](#job-header-document) | **REQUIRED**. |
| 1 | \[\][Job Document](#job-document) | **REQUIRED**. Define what the job to invoke. |

## Config Example

```yaml
---
type: web hook
---
- name: example_http_job
  cron:
    express: "* * * * *"
    skip_if_still_running: false
    delay_if_still_running: true
  config:
    uri: http://localhose/example
    parameters:
      - name: rid
        value: $request_id
```

## Schema

### Job Header Document

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| type | string | **REQUIRED**. The available values are '*web hook*' and '*file*'. |

### Job Document

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| name | string | The given name of the job. |
| cron | [Cron Config Object](#cron-config-object) | The cron setting. |
| config | [HTTP Job Config Object](#http-job-config-object) \| [File Job Config Object](#file-job-config-object) | **REQUIRED**. The type of the field depends on the value of the `type` field in [index-0 document](#job-header-document). |

### Cron Config Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| express | string | An expression represents a set of times, using 5 space-separated fields like '`* * * * *`'. The form is based on [the Cron](https://en.wikipedia.org/wiki/Cron). |
| skip_if_still_running | boolean | Skip the running job if true. |
| delay_if_still_running | boolean | Run the job util the job is done. |

### HTTP Job Config Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| uri | string | **REQUIRED**. A valid URI with http or https schema. |
| parameters | \[\][Job Config Value Object](#job-config-value-object) | Parameters are sent as a JSON-form request body. |

### File Job Config Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| path | string | **REQUIRED**. An executable file path. |
| flags | \[\][Job Config Value Object](#job-config-value-object) | Set the given flags if any. |

### Job Config Value Object

Fixed fields:

| Field Name | Type | Description |
| --- | --- | --- |
| name | string | **REQUIRED**. The given name of the parameter. |
| value | string | **REQUIRED**. See [Job Config Value List](#job-config-value-list) for details. |

### Job Config Value List

| value | comment |
| --- | --- |
| `$request_id` | A version 4 [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier), it is useful to track the scope. |
