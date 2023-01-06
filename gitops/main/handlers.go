package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/exp/maps"

	//"github.com/kr/pretty"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"

	"github.com/polarbroadband/rp1/etcdlib"
	"github.com/polarbroadband/rp1/gitlib"
	"github.com/polarbroadband/rp1/proto/dns"
	"gopkg.in/yaml.v3"
)

var (
	IaC_PATTERN_DNS  = `(?i)^DNS_.*?\.ya?ml$`
	IaC_PATTERN_DHCP = `(?i)^DHCP_.*?\.ya?ml$`
)

type MetaIaC struct {
	ApiVersion string `json:"ApiVersion" yaml:"apiVersion"`
	Kind       string `json:"Kind" yaml:"kind"`
	MetaData   struct {
		Name   string            `json:"Name" yaml:"name"`
		Labels map[string]string `json:"Labels" yaml:"labels"`
	} `json:"MetaData" yaml:"metadata"`
}

type DnsZone struct {
	Spec *dns.Zone `json:"Spec" yaml:"spec"`
}

// healtz response k8s health check probe
func (api *Gateway) healtz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "release": "unknown"}); err != nil {
		api.Error(w, 500, err.Error(), "server error")
	}
}

func (api *Gateway) WebhookV1Push(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if event := r.Header.Get("X-Gitlab-Event"); event != "Push Hook" {
		api.Error(w, http.StatusBadRequest, fmt.Sprintf("invalid webhook event %s, from %s", event, r.Header.Get("X-Gitlab-Instance")))
		return
	}

	if r.Header.Get("X-Gitlab-Token") != GITOPS_API_TOKEN {
		api.Error(w, http.StatusUnauthorized, fmt.Sprintf("invalid webhook token, request from %s", r.Header.Get("X-Gitlab-Instance")), "Unauthorized")
		return
	}

	var event struct {
		Commit string `json:"checkout_sha"`
		Repo   struct {
			ID   float64 `json:"id"`
			Name string  `json:"name"`
		} `json:"project"`
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		api.Error(w, http.StatusBadRequest, fmt.Sprintf("unable to parse webhook event %v", err))
		return
	}

	commit := gitlib.GitLabCommit{
		GitLab: api.GitLab,
		Repo:   event.Repo.ID,
		Commit: event.Commit,
		Log: api.Log.WithFields(logrus.Fields{
			"git":    r.Header.Get("X-Gitlab-Instance"),
			"repo":   event.Repo.Name,
			"commit": event.Commit,
		}),
	}

	commit.Log.Info("webhook push event")

	switch strings.ToLower(event.Repo.Name) {
	case "dns":
		blobs, err := commit.GetRepoRawFiles(regexp.MustCompile(IaC_PATTERN_DNS))
		if err != nil {
			api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to get content %v", err))
			return
		}
		commitData := &dns.Zone{
			Commit:  event.Commit,
			Records: make(map[string]*dns.Category),
		}
		for _, b := range blobs {
			commit.Log.Infof("processing file: %s", b.Path)

			var meta MetaIaC
			err = yaml.Unmarshal(b.Content, &meta)

			if err != nil {
				api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse meta of %s: %v", b.Path, err))
				return
			}

			//pretty.Printf("\n--- Meta STRUCT ---\n%# v\n\n", meta)

			if meta.Kind != "DNS" {
				api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse data of %s: invalid service kind", b.Path))
				return
			}
			if meta.ApiVersion == "service/v1" {
				// v1 DNS data model
				var data DnsZone
				err = yaml.Unmarshal(b.Content, &data)

				if err != nil {
					api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to parse data of %s: %v", b.Path, err))
					return
				}

				//pretty.Printf("\n--- DNS STRUCT ---\n%# v\n\n", data.Spec)
				maps.Copy(commitData.Records, data.Spec.Records)
			}
			// v2 ... model
		}

		out, err := proto.Marshal(commitData)
		if err != nil {
			api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to serialize data %v", err))
			return
		}
		depot := etcdlib.NewKvDepot(ETCD_IaC_DNS, api.Client, api.Log)
		if err = depot.Put("zone", string(out), 0); err != nil {
			api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to publish data %v", err))
			return
		}

	case "dhcp":
		blobs, err := commit.GetRepoRawFiles(regexp.MustCompile(IaC_PATTERN_DHCP))
		if err != nil {
			api.Error(w, http.StatusInternalServerError, fmt.Sprintf("unable to get content %v", err))
			return
		}
		//pretty.Printf("\n--- Selected BLOB STRUCT ---\n%# v\n\n", blobs)
		for _, b := range blobs {
			commit.Log.Infof("file: %s +++++++++ RAW: <<%v>>", b.Path, b.Content)
		}

	default:
		api.Error(w, http.StatusInternalServerError, fmt.Sprintf("invalid repo name"))
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		api.Error(w, 500, err.Error(), "server error")
	}
}
