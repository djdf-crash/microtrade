<h1>Schema config file</h1>
*******************************
 1 {
 2   "definitions": {}, 
 3   "$schema": "http://json-schema.org/draft-06/schema#", 
 4   "$id": "config.json", 
 5   "type": "object", 
 6   "additionalProperties": false, 
 7   "properties": {
 8     "mode_start": {
 9       "type": "string"
10     }, 
11     "port": {
12       "type": "string"
13     }, 
14     "send_email": {
15       "type": "object", 
16       "additionalProperties": false, 
17       "properties": {
18         "server": {
19           "type": "string"
20         }, 
21         "port": {
22           "type": "string"
23         }, 
24         "sender": {
25           "type": "string"
26         }, 
27         "password_sender": {
28           "type": "string"
29         }
30       }, 
31       "required": [
32         "server", 
33         "port", 
34         "sender", 
35         "password_sender"
36       ]
37     }, 
38     "data_base": {
39       "type": "object", 
40       "additionalProperties": false, 
41       "properties": {
42         "name_driver": {
43           "type": "string"
44         }, 
45         "path": {
46           "type": "string"
47         }
48       }, 
49       "required": [
50         "name_driver", 
51         "path"
52       ]
53     }
54   }, 
55   "required": [
56     "mode_start", 
57     "port", 
58     "send_email", 
59     "data_base"
60   ]
61 }