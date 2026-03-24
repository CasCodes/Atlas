# Atlas
A passive network scanner that visualizes your traffic as a live world map. Written in Go using gopacket and D3.js.

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
