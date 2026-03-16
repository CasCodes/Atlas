# GoScan

Run the api server with `sudo go run .`

## Endpoints:

### POST /scan
Runs the scan on env0 for the specified milisecond duration
```
query args: 
- <duration in ms> (int)

usage:
curl -X POST "http://localhost:8080/scan?duration=5000"
```

### GET /graph
Returns the json serialized graph of all captured packages 
```
query args:

usage:
curl "http://localhost:8080/graph"
```


## To do
Todo:
- specify scan duration with arg
- build a graph over the captures destinations, make make edges larger based on number of requests. 


Backlog:
- MCP server over graph