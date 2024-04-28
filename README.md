# golang-multicast

UDP multicast implemented in golang.

## Run it

Run as many instances of `go run broadcaster/broadcaster.go` and `go run receiver/receiver.go` as you want. It will default to use UDP-based multicast.

In order to use IP-multicast you need to specify the interface to use for multicasting by providing `--itfname eth0` (change to your desired interface) as a flag to the broadcaster and receiver.