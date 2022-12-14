package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/polarbroadband/rp1/etcdlib"
	"github.com/polarbroadband/rp1/proto/dns"
	etcd "go.etcd.io/etcd/client/v3"

	pkgdns "github.com/miekg/dns"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetReportCaller(true)
}

var (
	ETCD_IaC_DNS   = os.Getenv("ETCD_IaC_DNS") // "/cirrus/iac/dns"
	ETCD_ENDPOINTS = strings.Split(os.Getenv("ETCD_ENDPOINTS"), ",")
	ETCD_USERNAME  = os.Getenv("ETCD_USERNAME")
	ETCD_PASSWORD  = os.Getenv("ETCD_PASSWORD")
	READY          = false
)

type BaseDNS struct {
	Locker *sync.RWMutex
	*dns.Zone
	Log *logrus.Entry
}

func (b *BaseDNS) Update(data *dns.Zone) {
	b.Locker.Lock()
	defer b.Locker.Unlock()
	b.Zone = data
}

func (b *BaseDNS) Search(fqdn, cat string) *dns.Record {
	b.Locker.RLock()
	defer b.Locker.RUnlock()
	if c, ok := b.GetRecords()[fqdn]; ok {
		if v, exist := c.Type[cat]; exist {
			return v
		}
	}
	return nil
}

func (b *BaseDNS) String() (str string) {
	b.Locker.RLock()
	defer b.Locker.RUnlock()
	str = fmt.Sprintf("\nCommit: %s", b.Commit)
	for fqdn, c := range b.GetRecords() {
		for t, v := range c.Type {
			str = str + fmt.Sprintf("\nFQDN: %v --- Type: %v --- Addr: %+v --- TTL: %v", fqdn, t, v.GetAddr(), v.GetTTL())
		}
	}
	return
}

func (b *BaseDNS) ServeDNS(w pkgdns.ResponseWriter, r *pkgdns.Msg) {
	msg := pkgdns.Msg{}
	msg.SetReply(r)
	for _, q := range r.Question {
		log := b.Log.WithFields(logrus.Fields{"src": w.RemoteAddr().String(), "type": q.Qtype, "domain": q.Name})
		switch q.Qtype {
		case pkgdns.TypeA:
			msg.Authoritative = true
			domain := q.Name
			record := b.Search(domain, "A")
			if record != nil {
				log.Info("DNS query received, found records")
				ttl := uint32(60)
				if record.TTL != 0 {
					ttl = uint32(record.TTL)
				}
				for _, addr := range record.GetAddr() {
					msg.Answer = append(msg.Answer, &pkgdns.A{
						Hdr: pkgdns.RR_Header{Name: domain, Rrtype: pkgdns.TypeA, Class: pkgdns.ClassINET, Ttl: ttl},
						A:   net.ParseIP(addr),
					})
				}
			} else {
				log.Info("DNS query received, no records")
			}
		case pkgdns.TypeAAAA:
			msg.Authoritative = true
			domain := q.Name
			record := b.Search(domain, "AAAA")
			if record != nil {
				log.Info("DNS query received, found records")
				ttl := uint32(60)
				if record.TTL != 0 {
					ttl = uint32(record.TTL)
				}
				for _, addr := range record.GetAddr() {
					msg.Answer = append(msg.Answer, &pkgdns.AAAA{
						Hdr:  pkgdns.RR_Header{Name: domain, Rrtype: pkgdns.TypeAAAA, Class: pkgdns.ClassINET, Ttl: ttl},
						AAAA: net.ParseIP(addr),
					})
				}
			} else {
				log.Info("DNS query received, no records")
			}
		default:
			log.Info("DNS query received, unsupported type")
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		logrus.Fatal(err)
	}
	log := logrus.WithFields(logrus.Fields{"wkr": hostname, "pkg": "dns_server"})

	etcdClientCfg := etcd.Config{
		Endpoints:   ETCD_ENDPOINTS,
		Username:    ETCD_USERNAME,
		Password:    ETCD_PASSWORD,
		DialTimeout: etcdlib.DEFAULT_ETCD_DIAL_TIMEOUT,
	}
	if dialTimeout, err := time.ParseDuration(os.Getenv("ETCD_DIAL_TIMEOUT")); err != nil {
		log.Warnf("invalid env variable ETCD_DIAL_TIMEOUT: %v, set to %v", err, etcdClientCfg.DialTimeout)
	} else {
		etcdClientCfg.DialTimeout = dialTimeout
	}
	etcdClient, err := etcd.New(etcdClientCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer etcdClient.Close()

	depot := etcdlib.NewKvDepot(ETCD_IaC_DNS, etcdClient, log)
	current, zoneCH, err := depot.Subscribe("zone")
	if err != nil {
		log.Fatal(err)
	}
	defer depot.Cancel()

	base := BaseDNS{&sync.RWMutex{}, &dns.Zone{}, log}

	if current == nil {
		log.Warnf("server not ready, zone data not available")
	} else {
		if err := proto.Unmarshal(current.Value, base.Zone); err != nil {
			log.Fatal(err)
		}
		READY = true
		fmt.Println("Current " + base.String())
	}

	go func() {
		log.Infof("start zone watcher %s/zone", ETCD_IaC_DNS)
		for wresp := range zoneCH {
			ev := wresp.Events[0]
			if err := proto.Unmarshal(ev.Kv.Value, base.Zone); err != nil {
				log.Errorf("received invalid zone data %v", err)
				continue
			}
			fmt.Println("New " + base.String())

			READY = true
			fmt.Println("next loop")
		}
	}()

	srv := &pkgdns.Server{
		Addr:    ":53",
		Net:     "udp",
		Handler: &base,
	}
	log.Infof("start DNS listener")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %v\n", err)
	}
}
