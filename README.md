# gocron

Simple golang wrapper over `github.com/robfig/cron` and `os/exec` as a cron replacement.
Additionally the application opens a HTTP port that can be used as a healthcheck.2

## Getting Started

Install the dependencies:

```bash
go get .
```

Usage:

```bash
go run ./main.go -s "* * * * *" -p 8080 -i -- /bin/bash -c "echo 1;"
```

Check the healthcheck:

```
$ curl -v localhost:8080
* Rebuilt URL to: localhost:8080/
* Hostname was NOT found in DNS cache
*   Trying ::1...
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> User-Agent: curl/7.37.1
> Host: localhost:8080
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 11 Mar 2015 12:59:07 GMT
< Content-Length: 237
<
{
  "Running": {},
  "Status": {
    "Exit_status": 0,
    "Stdout": "1\n",
    "Stderr": "",
    "ExitTime": "2015-03-11T13:59:05+01:00",
    "Pid": 14420,
    "StartingTime": "2015-03-11T13:59:05+01:00"
  },
  "Schedule": "*/5 * * * *"
* Connection #0 to host localhost left intact
}
```

## License

This program is licensed under the MIT-License. See the "LICENSE" file for more information
