# kaf-proxy

## Run the example

```bash
go run . 
```

The server will be reachable at `http://localhost:8090`.

```bash
# True positive request (403 Forbidden)
curl -i 'localhost:8090/hello?id=0'
# True negative request (200 OK)
curl -i 'localhost:8090/hello'
```