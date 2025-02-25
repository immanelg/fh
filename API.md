API:
POST /list-dir
POST /read-file
POST /create-file
POST /create-dir
POST /delete
POST /copy
POST /change
POST /search
POST /log

API usage examples:

```
mkdir /tmp/test
echo contents > /tmp/test/myfile
go run . -d /tmp/test  # or air
```

```
curl -w ", %{http_code}\n" -X POST http://localhost:8080/create-file -F 'payload={"Path":"a","Type":"File"}' -F 'file=@/tmp/test/myfile'
curl -w ", %{http_code}\n" -X POST http://localhost:8080/create-dir -d '{"Path": "A"}'
curl -w ", %{http_code}\n" http://localhost:8080/read-file -d '{"Path": "a"}'
curl http://localhost:8080/list-dir -d '{"Path": "./"}' | jq
curl -w ", %{http_code}\n" http://localhost:8080/copy -d '{"Src": "a", "Dst": "v"}'
curl -w ", %{http_code}\n" http://localhost:8080/delete -d '{"Path": "a"}'
```
