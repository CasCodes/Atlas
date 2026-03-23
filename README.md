# Atlas
A small package capturing utility written in go. The captured packages are visualized as a graph, where nodes are IP adresses and edges are stronger the more packages were send between the nodes.

What it does:
```
Capture packages -> run traceroute on all discovered IPs -> plot the routes as a graph
```

## Usage
Run the api server with `sudo go run .`

Visit the web interface at `http://localhost:8080/`

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
- make the graph live, i.e. stream graph data continuosly (with start/stop scan)

Backlog:
- MCP server over graph