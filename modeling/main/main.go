package main

import (
	//"context"
	//"fmt"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	//clientv3 "go.etcd.io/etcd/client/v3"

	yaml "gopkg.in/yaml.v3"

	"github.com/kr/pretty"

	//pb "github.com/polarbroadband/rp1/proto/dns"
	"github.com/polarbroadband/rp1/modeling/pb"
)

var (
	YFILE         = "../yamls/dns_t01.yml"
	ETCD_USERNAME = "root"
	ETCD_PASSWORD = "WT6kZlJ24Y"
	//ETCD_ENDPOINTS    = []string{"172.17.0.4:2379"}
	ETCD_ENDPOINTS    = []string{"etcd:2379"}
	ETCD_DIAL_TIMEOUT = time.Second * 5
	ETCD_OPS_TIMEOUT  = time.Second * 1
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
	Spec *pb.Zone `json:"Spec" yaml:"spec"`
}

func main() {

	file, err := ioutil.ReadFile(YFILE)
	if err != nil {
		log.Fatal(err)
	}

	var meta DnsZone
	err = yaml.Unmarshal(file, &meta)

	if err != nil {
		log.Fatal(err)
	}

	pretty.Printf("\n--- Meta STRUCT ---\n%# v\n\n", meta.Spec.Records["t01.cirrus.io"].Type["AAAA"].Addr)
	//fmt.Printf("Type-%T, Value-%v\n", meta.Spec.Records["t01.cirrus.io"].Type["none"], meta.Spec.Records["t01.cirrus.io"].Type["none"])

	if v, ok := meta.Spec.Records["none"].Type["AAAA"]; ok {
		fmt.Println(v)
	} else {
		fmt.Println("nil")
	}

}

/*
func main() {

	hold := make(chan struct{})

	file, err := ioutil.ReadFile(YFILE)
	if err != nil {
		log.Fatal(err)
	}

	var meta MetaIaC
	err = yaml.Unmarshal(file, &meta)

	if err != nil {
		log.Fatal(err)
	}

	pretty.Printf("\n--- Meta STRUCT ---\n%# v\n\n", meta)

	if meta.ApiVersion == "service/v1" && meta.Kind == "DNS" {
		var data DnsZone
		err = yaml.Unmarshal(file, &data)

		if err != nil {
			log.Fatal(err)
		}

		pretty.Printf("\n--- DNS STRUCT ---\n%# v\n\n", data.Spec)

		out, err := proto.Marshal(data.Spec)
		if err != nil {
			log.Fatalln("Failed to encode data:", err)
		}

		cli, err := clientv3.New(clientv3.Config{
			Endpoints:   ETCD_ENDPOINTS,
			DialTimeout: ETCD_DIAL_TIMEOUT,
			Username:    ETCD_USERNAME,
			Password:    ETCD_PASSWORD,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer cli.Close()

		_, err = cli.Put(context.TODO(), "foo", "bar")
		if err != nil {
			log.Fatal(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), ETCD_OPS_TIMEOUT)
		resp, err := cli.Get(ctx, "foo")

		if err != nil {
			log.Fatal(err)
		}
		for _, ev := range resp.Kvs {
			fmt.Printf("%s : %s\n", ev.Key, ev.Value)
		}

		go func() {
			rch := cli.Watch(context.Background(), "dns")
			for wresp := range rch {
				for _, ev := range wresp.Events {
					fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					var re pb.Zone
					if err := proto.Unmarshal(ev.Kv.Value, &re); err != nil {
						log.Fatal(err)
					}
					re2 := re.GetRecords()[2]
					pretty.Printf("\n--- WATCHER-1 ---\n%# v\n\n", re2.GetTTL())
				}
			}
		}()

		go func() {
			rch := cli.Watch(context.Background(), "dns")
			for wresp := range rch {
				for _, ev := range wresp.Events {
					fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					var re pb.Zone
					if err := proto.Unmarshal(ev.Kv.Value, &re); err != nil {
						log.Fatal(err)
					}
					re2 := re.GetRecords()[2]
					pretty.Printf("\n--- WATCHER-2 ---\n%# v\n\n", re2.GetTTL())
				}
			}
		}()

		_, err = cli.Put(ctx, "dns", string(out))
		if err != nil {
			log.Fatal(err)
		}

		resp1, err := cli.Get(ctx, "dns")

		if err != nil {
			log.Fatal(err)
		}
		pretty.Printf("\n--- Retrieved RAW ---\n%# v\n\n", resp1)

		for _, ev := range resp1.Kvs {
			var re pb.Zone
			if err := proto.Unmarshal(ev.Value, &re); err != nil {
				log.Fatal(err)
			}
			pretty.Printf("\n--- Retrieved STRUCT ---\n%# v\n\n", re.GetRecords())
			//fmt.Printf("\n--- Retrieved STRUCT ---\n%+v", re)
		}

		go func() {
			for {
				if data.Spec.Records[2].TTL > 15 {
					break
				}
				data.Spec.Records[2].TTL++
				out, err := proto.Marshal(data.Spec)
				if err != nil {
					log.Fatalln("Failed to encode data:", err)
				}

				_, err = cli.Put(context.Background(), "dns", string(out))
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("\n+++++++++++++++++++++++ ADD TTL ++++\n")
				time.Sleep(time.Second)
			}

		}()

		<-hold
		cancel()

	} else {
		fmt.Println("invalid input")
	}

}
*/

/*
type KvRepo struct {
	Repo      *clientv3.Client
	Prefix    string
	OprCtx    context.Context
	OprCancel context.CancelFunc
}

func (kv *KvRepo) SetCTX() {
	kv.OprCtx, kv.OprCancel = context.WithTimeout(context.Background(), ETCD_OPS_TIMEOUT)
}

func (kv *KvRepo) HasPrefix(key string) (bool, error) {
	return false, nil
}

func (kv *KvRepo) FindLatest(key string) (*clientv3.Event, error) {
	kv.SetCTX()
	rch := kv.Repo.Watch(kv.OprCtx, kv.Prefix+"/"+key, clientv3.WithPrefix(), clientv3.WithRev(1)) //, clientv3.WithSort(clientv3.SortByModRevision, clientv3.SortDescend)
	for wresp := range rch {
		for i := len(wresp.Events) - 1; i >= 0; i-- {
			if wresp.Events[i].Kv.Version != 0 {
				kv.OprCancel()
				return wresp.Events[i], nil
			}
		}
		return nil, fmt.Errorf("invalid etcd watch event")
	}
	fmt.Println("++++++++++++ None +++++++++++++++")
	return nil, nil
}

func main() {

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   ETCD_ENDPOINTS,
		DialTimeout: ETCD_DIAL_TIMEOUT,
		Username:    ETCD_USERNAME,
		Password:    ETCD_PASSWORD,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	t01 := KvRepo{cli, "/lease", nil, nil}

	// t01.SetCTX()
	// options := append(clientv3.WithLastCreate(), clientv3.WithPrefix())
	// resp, err := t01.Repo.Get(t01.OprCtx, t01.Prefix+"/a", options...)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// pretty.Printf("\n--- RESP STRUCT ---\n%# v\n\n", resp)

	mn, _ := t01.FindLatest("a/a01")
	pretty.Printf("\n--- latest ---\n%# v\n\n", mn)

}
*/
