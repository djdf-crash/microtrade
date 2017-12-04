<h1>Schema config file</h1>
### JSON Schema

Here are JSON Schema.

```json
{
  "definitions": {},
  "$schema": "http://json-schema.org/draft-06/schema#",
  "$id": "config.json",
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "mode_start": {
      "type": "string"
    },
    "port": {
      "type": "string"
    },
    "send_email": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "server": {
          "type": "string"
        },
        "port": {
          "type": "string"
        },
        "sender": {
          "type": "string"
        },
        "password_sender": {
          "type": "string"
        }
      },
      "required": [
        "server",
        "port",
        "sender",
        "password_sender"
      ]
    },
    "data_base": {
      "type": "object",
      "additionalProperties": false,
      "properties": {
        "name_driver": {
          "type": "string"
        },
        "path": {
          "type": "string"
        }
      },
      "required": [
        "name_driver",
        "path"
      ]
    }
  },
  "required": [
    "mode_start",
    "port",
    "send_email",
    "data_base"
  ]
}
```

## Examples

```json

{
  "mode_start": "release", //start mode "debug" or "release"
  "port": ":8080", // start on port
  "send_email":{
      "server":"smtp.gmail.com", //server smtp
      "port":":587", //port smtp
      "sender":"user email sender",
      "password_sender":"password"
  },
  "data_base":{
      "name_driver":"sqlite3", //driver DB "sqlite3" or "MySql" or etc.
      "path":"path" //path driver sqlite3 "./dadabase.db" or MySql "user:password@/dbname"
    }
}

```
