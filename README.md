# tinycron
A very small replacement for cron

## Installing

```bash
curl -sLo https://github.com/bcicen/tinycron/releases/download/v0.1/tinycron-0.1-linux-amd64
sudo mv tinycron /usr/local/bin/
sudo chmod +x /usr/local/bin/tinycron
```

# Usage

```
tinycron [expression] [command]
```

Tinycron can be invoked via commandline:
```bash
$ /usr/local/bin/tinycron '*/5 * * * * * *' /bin/echo hello
```

Or used in your scripts interpreter line:
```bash
#!/usr/local/bin/tinycron */5 * * * * * * /bin/sh # run this script every five seconds
echo "Current time: $(date)"
```

## Expressions

Tinycron uses and supports all expressions in the [cronexpr](https://github.com/gorhill/cronexpr) library. Some examples:

`@daily` - run once daily, at midnight
`*/30 * * * * * *` - run every 30 seconds
`* 15 * * * * *` - run at minute `:15` of every hour

## Config

TinyCron can be configured by setting the below environmental variables to a non-empty value:

Variable | Description
--- | ---
TINYCRON_DEBUG | Enable debug output
