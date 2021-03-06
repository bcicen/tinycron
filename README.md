# tinycron
A very small replacement for cron. Particularly useful within containers and for distributing cron tasks alongside a project as a standalone file.

## Installing

```bash
curl -sLo tinycron https://github.com/bcicen/tinycron/releases/download/v0.3/tinycron-0.3-linux-amd64
sudo mv tinycron /usr/local/bin/
sudo chmod +x /usr/local/bin/tinycron
```

# Usage

```
tinycron [expression] [command]
```

Tinycron can be conveniently used in your scripts interpreter line:
```bash
#!/usr/local/bin/tinycron */5 * * * * * * /bin/sh
echo "Current time: $(date)"
```

Or invoked directly via commandline:
```bash
$ tinycron '*/5 * * * * * *' /bin/echo hello
```


## Expressions

Tinycron uses and supports expressions from the [cronexpr](https://github.com/gorhill/cronexpr) library. Some examples:

* `@daily` - run once daily, at midnight
* `* 15 * * * * *` - run at minute `:15` of every hour
* `*/30 * * * * * *` - run every 30 seconds

## Config

TinyCron can be configured by setting the below environmental variables to a non-empty value:

Variable | Description
--- | ---
TINYCRON_VERBOSE | Enable verbose output
