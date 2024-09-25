# halon-smtpd-exporter

Basic Prometheus exporter for Halon SMTPD server.

It uses the [control socket API](https://docs.halon.io/manual/api_control_sockets.html) to retrieve statistics from the local SMTP servers.

These are then formatted in the metrics format and served at `http//localhost:9393/metrics` (by default).

## Config

The following environment variables are available:

| Variable                         | Description                                                       |
| -------------------------------- | ----------------------------------------------------------------- |
| HALON_SMTPD_EXPORTER_SOCKET_PATH | Path to the control socket, default is `/var/run/halon/smtpd.ctl` |
| HALON_SMTPD_EXPORTER_LISTENADDR  | IP:Port for the webserver, default is `:9393`                     |

## Protobuf

You will need the Protobuf definitions matching your `smtpd` version, they can be found [here](https://docs.halon.io/protobuf-schemas/).

Clone this repo, grab the `smtpd.proto` file for your version.

Install `protoc-gen-go`:

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

Install `protoc`, for Debian/Ubuntu:

```
sudo apt install protobuf-compiler
```

Generate Golang types:

```
protoc -I=$(pwd) --go_out=$(pwd) --go_opt=Msmtpd.proto=pkg/halon_smtpd_ctl smtpd.proto
```
