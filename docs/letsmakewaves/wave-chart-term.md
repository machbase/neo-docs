---
title: Chart on terminal
parent: Let's make waves
layout: default
order: 30
nav_order: 210
---

# Chart on terminal

machabse-neo shell provides simple command line tool for monitoring incoming data.

```sh
machbase-neo shell chart \
    --range 30s \
    EXAMPLE/wave.sin#value EXAMPLE/wave.cos#value
```

![img](../img/term-chart02.gif)

