# fh
File Storage API over HTTP.

# Running
## Dev server
I recommend to use `just` and `watchexec`:
```sh
just
```
## Example API requests
```sh
curl -X GET -w "%{http_code}\n" http://localhost:8080/one
curl -X GET -w "%{http_code}\n" http://localhost:8080/
curl -X GET -w "%{http_code}\n" http://localhost:8080/one?metadata=1
curl -X POST -w "%{http_code}\n" --data-binary "@tmp/one" http://localhost:8080/two
curl -X POST -w "%{http_code}\n" http://localhost:8080/newdir?resource-type=dir
curl -X POST -w "%{http_code}\n" -H "Source-Path: one" -H "Operation: Copy" http://localhost:8080/two
curl -X DELETE -w "%{http_code}\n" http://localhost:8080/two
```

# See also
* [WebDAV](https://en.wikipedia.org/wiki/WebDAV)
* [9P](https://en.wikipedia.org/wiki/9P_(protocol))
* [sshfs](https://github.com/libfuse/sshfs)
