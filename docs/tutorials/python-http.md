---
layout: default
title: HTTP API in Python
parent: Tutorials
has_children: false
nav_order: 91
---

# HTTP API in Python
{: .no_toc }

1. TOC
{: toc }

## Query

### GET CSV

```python
import requests
params = {"q":"select * from example", "format":"csv", "heading":"false"} 
response = requests.get("http://127.0.0.1:5654/db/query", params)
print(response.text)
```

## Write

### POST CSV

```python
import requests
csvdata = """temperature,1677033057000000000,21.1
humidity,1677033057000000000,0.53
"""
response = requests.post(
    "http://127.0.0.1:5654/db/write/example?heading=false", 
    data=csvdata, 
    headers={'Content-Type': 'text/csv'})
print(response.json())
```
