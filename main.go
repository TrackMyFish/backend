package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/trackmyfish/backend/internal/server"
	trackmyfishv1alpha1 "github.com/trackmyfish/proto/trackmyfish/v1alpha1"
)

//go:embed build
var feStatic embed.FS

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	logrus.SetOutput(os.Stdout)

	// Only log the info severity or above.
	logrus.SetLevel(logrus.InfoLevel)
}

func handleBindEnvErr(err error) {
	if err != nil {
		logrus.Fatalf("unable to bind viper key to environment variable: '%+v'", err)
	}
}

func main() {
	// Config files
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Custom config file mapped as a volume when using Docker
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/config")

	// Environment Variables
	handleBindEnvErr(viper.BindEnv("server.port", "TMF_SERVER_PORT"))
	handleBindEnvErr(viper.BindEnv("server.httpProxy.enabled", "TMF_HTTP_PROXY_ENABLED"))
	handleBindEnvErr(viper.BindEnv("server.httpProxy.port", "TMF_HTTP_PROXY_PORT"))
	handleBindEnvErr(viper.BindEnv("db.host", "TMF_DB_HOST"))
	handleBindEnvErr(viper.BindEnv("db.port", "TMF_DB_PORT"))
	handleBindEnvErr(viper.BindEnv("db.username", "TMF_DB_USERNAME"))
	handleBindEnvErr(viper.BindEnv("db.password", "TMF_DB_PASSWORD"))
	handleBindEnvErr(viper.BindEnv("db.name", "TMF_DB_NAME"))

	// Merge config
	if err := viper.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore as we use defaults/environment variables
			// and if anything required isn't set (e.g. db password) we'll error later on
		} else {
			// Config file was found but another error was produced
			logrus.Fatalf("unable to merge in config: '%+v'", err)
		}
	}

	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.httpProxy.enabled", false)
	viper.SetDefault("server.httpProxy.port", 8443)

	// DB defaults
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.username", "trackmyfish")
	viper.SetDefault("db.password", "")
	viper.SetDefault("db.name", "trackmyfish")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore as we use defaults/environment variables
			// and if anything required isn't set (e.g. db password) we'll error later on
		} else {
			// Config file was found but another error was produced
			logrus.Fatal("unable to read config: ", err)
		}
	}

	var (
		port             = viper.GetInt("server.port")
		httpProxyEnabled = viper.GetBool("server.httpProxy.enabled")
		httpProxyPort    = viper.GetInt("server.httpProxy.port")

		dbHost     = viper.GetString("db.host")
		dbPort     = viper.GetString("db.port")
		dbUsername = viper.GetString("db.username")
		dbPassword = viper.GetString("db.password")
		dbName     = viper.GetString("db.name")
	)

	logrus.WithFields(logrus.Fields{
		"Server Port":        port,
		"HTTP Proxy Enabled": httpProxyEnabled,
		"HTTP Proxy Port":    httpProxyPort,
		"Database Name":      dbName,
		"Database Host":      dbHost,
		"Database Port":      dbPort,
		"Database Username":  dbUsername,
	}).Info("Config Initialised")

	server, err := server.New(
		server.Config{DBHost: dbHost, DBPort: dbPort, DBUsername: dbUsername, DBPassword: dbPassword, DBName: dbName},
	)
	if err != nil {
		logrus.Fatalf("Unable to initialise new Server: %+v", err)
	}

	gServer := grpc.NewServer()

	trackmyfishv1alpha1.RegisterTrackMyFishServiceServer(gServer, server)

	reflection.Register(gServer)

	addr := fmt.Sprintf(":%d", port)

	if httpProxyEnabled {
		go httpProxyServer(httpProxyPort, addr)
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err, "Failed to create listener")
	}

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting grpc server")

	if err := gServer.Serve(listener); err != nil {
		logrus.Fatal(err, "Failed to start server")
	}
}

// httpProxyServer starts a new http server listening on the specified port, proxying
// requests to the provided grpc service
func httpProxyServer(port int, grpcAddr string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	grpcMux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := trackmyfishv1alpha1.RegisterTrackMyFishServiceHandlerFromEndpoint(ctx, grpcMux, grpcAddr, opts); err != nil {
		logrus.Fatal(err, "Failed to register http handler")
	}

	r := http.NewServeMux()

	r.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// gateway is generated to match for /v1alpha1/ and not /api/v1alpha1
		// we could update the gateway proto to match for /api/v1alpha1 but
		// it shouldn't care where it's mounted to, hence we just rewrite the path here
		r.URL.Path = strings.Replace(r.URL.Path, "/api", "", -1)
		grpcMux.ServeHTTP(w, r)
	})

	sch, err := buildHandler()
	if err != nil {
		logrus.Fatal(err, "unable to initialize build handler")
	}
	r.Handle("/", sch)

	logrus.WithFields(logrus.Fields{
		"port": port,
	}).Info("Starting http proxy server")

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), r), "Failed to start http proxy server")
}

func buildHandler() (http.Handler, error) {
	fsys := fs.FS(feStatic)

	sc, err := fs.Sub(fsys, "build")
	if err != nil {
		return nil, err
	}

	return http.FileServer(http.FS(sc)), nil
}
