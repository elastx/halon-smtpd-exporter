package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elastx/halon-smtpd-exporter/pkg/halon_smtpd_ctl"
	"github.com/golang/protobuf/proto"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	socketPath := os.Getenv("HALON_SMTPD_EXPORTER_SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/var/run/halon/smtpd.ctl"
	}

	h := halon_smtpd_ctl.New(socketPath)

	data, err := h.Query("q")
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	var ps halon_smtpd_ctl.ProcessStatsResponse
	err = proto.Unmarshal(data, &ps)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	output, err := ps.MarshalPrometheus()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "%s", err)
		return
	}

	fmt.Fprintf(w, "%s", output)
}

func main() {
	listenAddr := os.Getenv("HALON_SMTPD_EXPORTER_LISTENADDR")
	if listenAddr == "" {
		listenAddr = ":9393"
	}

	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(listenAddr, nil)
}
