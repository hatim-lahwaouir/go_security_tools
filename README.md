# tcp_port_scanner
A TCP port scanner is a small Go script that scans open TCP ports, powered by Go's goroutines and channels for fast scanning.

## Usage

```python
./main --help
Usage of ./main:
  -host string
        host to scan (default "locolhost")
  -last_port uint
        last port Number (default 3000)
  -rps uint
        request per second (default 100)
  -start_port uint
        starting port Number (default 1)
```


## Example

```python
./main -host google.com -start_port 400 -last_port 500
start scanning : google.com
port is up 443
```