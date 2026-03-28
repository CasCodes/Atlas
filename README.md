# Atlas
A passive network scanner that visualizes your traffic as a live world map. Written in Go using gopacket and D3.js.
Atlas scans for local traffic and runs traceroute & whois on all captured IPs along the path, then stream the resulting graph to the client.

<img width="1365" height="746" alt="atlas_example" src="https://github.com/user-attachments/assets/f558eb20-6a74-4ea0-9ace-a197ed51679c" />

> Larger edges mean more packages send on this path

## Usage
Run the api server with `sudo go run .`
> Sudo is required by some systems to read from the network interface

Visit the web interface at `http://localhost:8080/`

## Endpoints:

### POST /scanDuration
Runs the scan (synchronous) on env0 for the specified milisecond duration
```
query args: 
- duration (int, in miliseconds)

usage:
curl -X POST "http://localhost:8080/scan?duration=5000"
```

### GET /graph
Returns the json serialized graph of all captured packages 
```
Content-Type: application/json

usage:
curl "http://localhost:8080/graph"

response format:
{
  "nodes": [
    {
      "ip":      "51.104.15.252",
      "country": "Netherlands",
      "city":    "Amsterdam",
      "org":     "Microsoft Corporation",
      "lat":     52.3740,
      "lon":     4.8897
    }
  ],
  "edges": [
    {
      "from":  "192.168.178.191",
      "to":    "51.104.15.252",
      "count": 42
    }
  ]
}
```

### POST /start
Start running the scanner without timeout duration
```
query args: none

usage:
curl "http://localhost:8080/start"
```

### POST /stop
Stop any active scan
```
query args: none

usage:
curl "http://localhost:8080/stop"
```


### GET /stream
Opens a persistent SSE connection. The server pushes the full graph state every second.
The connection stays open until the client disconnects or the scan is stopped.

```
Content-Type: text/event-stream

usage:
curl -N "http://localhost:8080/stream"

response format:
data: {"nodes":[...],"edges":[...]}\n\n
```

