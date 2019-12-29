## Getting Started TURN

`turnserver -c /usr/local/etc/turnserver.conf --user user1:pass1 -X 0.0.0.0/10.146.0.4`

### TURN reader

```
TURN_SERVER=34.84.50.204
go run reader/main.go -h ${TURN_SERVER} -p 9610
```
