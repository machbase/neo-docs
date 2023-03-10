---
title: Read CSV with Pandas
parent: Tutorials
layout: default
nav_order: 5
---

# Pandas DataFrame directly from HTTP API.

This example shows how to load data into pandas dataframe via machbase-neo HTTP API.

![img](img/python_http_csv.jpg)

## Load CSV

Import pandas and urllib.

```py
from urllib import parse
import pandas as pd
```

Make query url for `"format": "csv"` option, then call `read_csv`.
Use `timeformat` to specify the precision of time data. `s`, `ms`, `us` and `ns`(default) are available.

```py
query_param = parse.urlencode({
    "q":"select * from example order by time limit 500",
    "format": "csv",
    "timeformat": "s",
})
df = pd.read_csv(f"http://127.0.0.1:5654/db/query?{query_param}")
df
```

## Load Compressed CSV

Read gzip'ed CSV from HTTP API.

```py
from urllib import parse
import pandas as pd
import requests
import io
```

```py
query_param = parse.urlencode({
    "q":"select * from example order by time desc limit 1000",
    "format": "csv",
    "timeformat": "s",
    "compress": "gzip",
})
response = requests.get(f"http://127.0.0.1:5654/db/query?{query_param}", timeout=30, stream=True)
df = pd.read_csv(io.BytesIO(response.content))
df
```

