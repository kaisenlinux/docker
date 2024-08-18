package main

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/containerd/log"
	metrics "github.com/docker/go-metrics"
)

func startMetricsServer(addr string) error {
	if addr == "" {
		return nil
	}
	if err := allocateDaemonPort(addr); err != nil {
		return err
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.Handler())
	go func() {
<<<<<<< HEAD
		log.G(context.TODO()).Infof("metrics API listening on %s", l.Addr())
		srv := &http.Server{
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Minute, // "G112: Potential Slowloris Attack (gosec)"; not a real concern for our use, so setting a long timeout.
		}
		if err := srv.Serve(l); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.G(context.TODO()).WithError(err).Error("error serving metrics API")
=======
		logrus.Infof("metrics API listening on %s", l.Addr())
		if err := http.Serve(l, mux); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			logrus.WithError(err).Error("error serving metrics API")
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		}
	}()
	return nil
}
