package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	proxy "github.com/kloudlite/kl/domain/dev-proxy"
	fn "github.com/kloudlite/kl/pkg/functions"
)

type Server struct {
	bin string
}

func New(binName string) *Server {
	return &Server{
		bin: binName,
	}
}
func portAvailable(port string) bool {
	address := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return false
	}
	defer listener.Close()
	return true
}

func (s *Server) Start(ctx context.Context) error {

	ch := make(chan error)

	defer func() {
		ctx.Done()
	}()

	app := http.NewServeMux()
	app.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		outputCh := make(chan string)
		errCh := make(chan error)

		command := strings.TrimPrefix(req.URL.Path, "/")

		switch command {
		case "healthy":
			w.WriteHeader(http.StatusOK)
			return

		case "exit":
			w.WriteHeader(http.StatusOK)
			ch <- nil
			return

		case "start", "stop", "status", "restart":

			go fn.StreamOutput(fmt.Sprintf("%s vpn %s", s.bin, command), map[string]string{"KL_APP": "true"}, outputCh, errCh)

			for {
				select {
				case output := <-outputCh:
					w.Write([]byte(output))
					w.(http.Flusher).Flush()
				case err := <-errCh:
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
					return
				}
			}
		default:
			w.WriteHeader(http.StatusNotFound)
			return
		}
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", proxy.AppPort),
		Handler: app,
	}

	fn.Logf("starting server at :%d", proxy.AppPort)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			ch <- err
		}
	}()

	err := <-ch

	if err2 := server.Shutdown(ctx); err2 != nil {
		return err2
	}

	return err
}
