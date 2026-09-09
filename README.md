# Cookie Test

This repo contains a small, self-contained Golang executable that demonstrates whether path scoping cookies prevents
them from being used across path scopes.

## URLs

* GET /debug - an HTML page that debugs the cookies available under the root path
* GET /app/debug - an HTML page that debugs the cookies available under the /app path
* GET /echo - a text/plain response echo'ing the cookies present in the request
* GET /app/echo - a text/plain response echo'ing the cookies present in the request
* GET /app/setcookies - sets various cookies used for testing various scenarios (Path, SameSite, etc)
* GET /app/clearcookies - clears any cookies that were set

## Running

You will need Golang installed. Once installed, the following will start the server:

```shell
go run main.go
```

At this point, the server is listening on port 10101. [Set your cookies](http://localhost:10101/app/setcookies)
and then view the [root debug](http://localhost:10101/debug) and [/app debug](http://localhost:10101/app/debug).
