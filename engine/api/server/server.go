package server // import "github.com/docker/docker/api/server"

import (
	"context"
	"net/http"
<<<<<<< HEAD
=======
	"strings"
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)

	"github.com/containerd/log"
	"github.com/docker/docker/api/server/httpstatus"
	"github.com/docker/docker/api/server/httputils"
	"github.com/docker/docker/api/server/middleware"
	"github.com/docker/docker/api/server/router"
	"github.com/docker/docker/api/server/router/debug"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/dockerversion"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// versionMatcher defines a variable matcher to be parsed by the router
// when a request is about to be served.
const versionMatcher = "/v{version:[0-9.]+}"

// Server contains instance details for the server
type Server struct {
	middlewares []middleware.Middleware
}

// UseMiddleware appends a new middleware to the request chain.
// This needs to be called before the API routes are configured.
func (s *Server) UseMiddleware(m middleware.Middleware) {
	s.middlewares = append(s.middlewares, m)
}

<<<<<<< HEAD
func (s *Server) makeHTTPHandler(handler httputils.APIFunc, operation string) http.HandlerFunc {
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
=======
// Accept sets a listener the server accepts connections into.
func (s *Server) Accept(addr string, listeners ...net.Listener) {
	for _, listener := range listeners {
		httpServer := &HTTPServer{
			srv: &http.Server{
				Addr: addr,
			},
			l: listener,
		}
		s.servers = append(s.servers, httpServer)
	}
}

// Close closes servers and thus stop receiving requests
func (s *Server) Close() {
	for _, srv := range s.servers {
		if err := srv.Close(); err != nil {
			logrus.Error(err)
		}
	}
}

// serveAPI loops through all initialized servers and spawns goroutine
// with Serve method for each. It sets createMux() as Handler also.
func (s *Server) serveAPI() error {
	var chErrors = make(chan error, len(s.servers))
	for _, srv := range s.servers {
		srv.srv.Handler = s.createMux()
		go func(srv *HTTPServer) {
			var err error
			logrus.Infof("API listen on %s", srv.l.Addr())
			if err = srv.Serve(); err != nil && strings.Contains(err.Error(), "use of closed network connection") {
				err = nil
			}
			chErrors <- err
		}(srv)
	}

	for range s.servers {
		err := <-chErrors
		if err != nil {
			return err
		}
	}
	return nil
}

// HTTPServer contains an instance of http server and the listener.
// srv *http.Server, contains configuration to create an http server and a mux router with all api end points.
// l   net.Listener, is a TCP or Socket listener that dispatches incoming request to the router.
type HTTPServer struct {
	srv *http.Server
	l   net.Listener
}

// Serve starts listening for inbound requests.
func (s *HTTPServer) Serve() error {
	return s.srv.Serve(s.l)
}

// Close closes the HTTPServer from listening for the inbound requests.
func (s *HTTPServer) Close() error {
	return s.l.Close()
}

func (s *Server) makeHTTPHandler(handler httputils.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
>>>>>>> parent of ea55db5 (Import the 20.10.24 version)
		// Define the context that we'll pass around to share info
		// like the docker-request-id.
		//
		// The 'context' will be used for global data that should
		// apply to all requests. Data that is specific to the
		// immediate function being called should still be passed
		// as 'args' on the function call.

		// use intermediate variable to prevent "should not use basic type
		// string as key in context.WithValue" golint errors
		ctx := context.WithValue(r.Context(), dockerversion.UAStringKey{}, r.Header.Get("User-Agent"))

		r = r.WithContext(ctx)
		handlerFunc := s.handlerWithGlobalMiddlewares(handler)

		vars := mux.Vars(r)
		if vars == nil {
			vars = make(map[string]string)
		}

		if err := handlerFunc(ctx, w, r, vars); err != nil {
			statusCode := httpstatus.FromError(err)
			if statusCode >= 500 {
				log.G(ctx).Errorf("Handler for %s %s returned error: %v", r.Method, r.URL.Path, err)
			}
			_ = httputils.WriteJSON(w, statusCode, &types.ErrorResponse{
				Message: err.Error(),
			})
		}
	}), operation).ServeHTTP
}

// CreateMux returns a new mux with all the routers registered.
func (s *Server) CreateMux(routers ...router.Router) *mux.Router {
	m := mux.NewRouter()

	log.G(context.TODO()).Debug("Registering routers")
	for _, apiRouter := range routers {
		for _, r := range apiRouter.Routes() {
			f := s.makeHTTPHandler(r.Handler(), r.Method()+" "+r.Path())

			log.G(context.TODO()).Debugf("Registering %s, %s", r.Method(), r.Path())
			m.Path(versionMatcher + r.Path()).Methods(r.Method()).Handler(f)
			m.Path(r.Path()).Methods(r.Method()).Handler(f)
		}
	}

	debugRouter := debug.NewRouter()
	for _, r := range debugRouter.Routes() {
		f := s.makeHTTPHandler(r.Handler(), r.Method()+" "+r.Path())
		m.Path("/debug" + r.Path()).Handler(f)
	}

	notFoundHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = httputils.WriteJSON(w, http.StatusNotFound, &types.ErrorResponse{
			Message: "page not found",
		})
	})

	m.HandleFunc(versionMatcher+"/{path:.*}", notFoundHandler)
	m.NotFoundHandler = notFoundHandler
	m.MethodNotAllowedHandler = notFoundHandler

	return m
}
