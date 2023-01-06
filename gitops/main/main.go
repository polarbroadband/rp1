package main

import (
	"crypto/tls"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/polarbroadband/rp1/etcdlib"
	"github.com/polarbroadband/rp1/gitlib"
	etcd "go.etcd.io/etcd/client/v3"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	//"github.com/kr/pretty"
	"github.com/polarbroadband/goto/util"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
}

var (
	GITOPS_API_TOKEN = os.Getenv("GITOPS_API_TOKEN")
	GIT_API_URL      = os.Getenv("GIT_API_URL")
	GIT_ACCESS_TOKEN = os.Getenv("GIT_ACCESS_TOKEN")
	ETCD_ENDPOINTS   = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	ETCD_USERNAME    = os.Getenv("ETCD_USERNAME")
	ETCD_PASSWORD    = os.Getenv("ETCD_PASSWORD")
	ETCD_IaC_DNS     = os.Getenv("ETCD_IaC_DNS")  // "/cirrus/iac/dns"
	ETCD_IaC_DHCP    = os.Getenv("ETCD_IaC_DHCP") // "/cirrus/iac/dhcp"
)

// Gateway holds arguments for calling api handler functions
type Gateway struct {
	*util.API
	*gitlib.GitLab
	*etcd.Client
}

func main() {
	hold := make(chan struct{})

	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatal(err)
	}

	etcdClientCfg := etcd.Config{
		Endpoints:   ETCD_ENDPOINTS,
		Username:    ETCD_USERNAME,
		Password:    ETCD_PASSWORD,
		DialTimeout: etcdlib.DEFAULT_ETCD_DIAL_TIMEOUT,
	}
	if dialTimeout, err := time.ParseDuration(os.Getenv("ETCD_DIAL_TIMEOUT")); err != nil {
		logrus.Warnf("invalid env variable ETCD_DIAL_TIMEOUT: %v, set to %v", err, etcdClientCfg.DialTimeout)
	} else {
		etcdClientCfg.DialTimeout = dialTimeout
	}
	etcdClient, err := etcd.New(etcdClientCfg)
	if err != nil {
		logrus.Fatal(err)
	}
	defer etcdClient.Close()

	api := Gateway{
		API: &util.API{
			//TokenSec: TOKENSEC,
			NoAuth: []string{}, // exclude gRPC Login from JWT check
			Log:    logrus.WithFields(logrus.Fields{"wkr": hostname, "pkg": "gitops"}),
		},
		GitLab: &gitlib.GitLab{
			URL:         GIT_API_URL,
			AccessToken: GIT_ACCESS_TOKEN,
			Client: &http.Client{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				},
			},
		},
		Client: etcdClient,
	}

	// setup and run http server
	r := mux.NewRouter()
	rru := r.PathPrefix("/cirrus/iac").Subrouter()
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{`*`}),
		//handlers.AllowedHeaders([]string{"content-type", "X-Csrf-Token", "withcredentials", "credentials",}),
		handlers.AllowedHeaders([]string{"content-type"}),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD"}),
	)

	rru.HandleFunc("/healtz", api.healtz).Methods("GET")

	rru.HandleFunc("/webhook/v1/push", api.WebhookV1Push).Methods("POST")

	go func() {
		api.Log.Info("GITOPS api server start")
		api.Log.Fatal(http.ListenAndServe(":8060", cors(rru)))
	}()

	<-hold
}
