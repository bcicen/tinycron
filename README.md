# tinycron
A very small replacement for cron

# Usage

TinyCron can be used directly in your scripts shebang line:
```
#!/usr/local/bin/tinycron */5 * * * * * * /bin/sh
echo "Current time: $(date)"
```

Or invoked via commandline:
```bash
$ /usr/local/bin/tinycron --debug '*/5 * * * * * * /bin/echo' 'hello!'
[tinycron] next job scheduled for 2016-05-16 11:58:45 -0400 EDT
[tinycron] running job: /bin/echo hello!
hello!
```
