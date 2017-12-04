<h1>Schema config file</h1>
### JSON Schema

Here are some examples to give you an idea how to use the class.

Assume you want to create the JSON object

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
