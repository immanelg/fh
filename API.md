API:
/create
/delete
/move
/copy
/search
/chmod
/list
/log

API usage examples:

```
mkdir /tmp/test

echo aaa > /tmp/test/a
echo bbb > /tmp/test/b

go run . -d /tmp/test
```

```
curl -w ", %{http_code}\n" http://localhost:8080/read -d '{"Path": "a"}'
curl -w ", %{http_code}\n" http://localhost:8080/copy -d '{"Src": "a", "Dst": "v"}'
curl -w ", %{http_code}\n" -X POST http://localhost:8080/create -F 'payload={"Path":"new","Type":"File"}' -F 'file=@/tmp/test/a'
curl http://localhost:8080/list -d '{"Path": "./"}' | jq
```
