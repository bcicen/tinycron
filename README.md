# tinycron
A very small replacement for cron

## Installing

```bash
curl -sLo https://github.com/bcicen/tinycron/releases/download/v0.1/tinycron-0.1-linux-amd64
sudo mv tinycron /usr/local/bin/
sudo chmod +x /usr/local/bin/tinycron
```

# Usage

TinyCron can be used directly in your scripts shebang line:
```
#!/usr/local/bin/tinycron */5 * * * * * * /bin/sh
echo "Current time: $(date)"
```

Or invoked via commandline:
```bash
$ /usr/local/bin/tinycron --debug '*/5 * * * * * *' /bin/echo hello
[tinycron] next job scheduled for 2016-05-16 11:58:45 -0400 EDT
[tinycron] running job: /bin/echo hello
hello
```

## Options

Option | Description
--- | ---
--debug, -d | Enable debug output
