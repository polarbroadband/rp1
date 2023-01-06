package etcdlib

import (
	"context"
	"fmt"
	"os"
	"time"

	etcd "go.etcd.io/etcd/client/v3"

	"github.com/sirupsen/logrus"
)

var (
	DEFAULT_ETCD_DIAL_TIMEOUT = time.Second * 5
	DEFAULT_ETCD_OPS_TIMEOUT  = time.Second * 3
)

type KvDepot struct {
	Depot      string
	Conn       *etcd.Client
	OprTimeout time.Duration
	Log        *logrus.Entry
	Cancel     context.CancelFunc
}

type KV struct {
	Key      string
	Value    []byte
	Lease    int64
	Revision int64
}

func NewKvDepot(depot string, conn *etcd.Client, log *logrus.Entry) *KvDepot {
	kv := KvDepot{
		Depot:      depot,
		Conn:       conn,
		OprTimeout: DEFAULT_ETCD_OPS_TIMEOUT,
		Log:        log.WithFields(logrus.Fields{"etcd": conn.Endpoints(), "depot": depot}),
	}
	if opsTimeout, err := time.ParseDuration(os.Getenv("ETCD_OPS_TIMEOUT")); err != nil {
		kv.Log.Warnf("invalid env variable ETCD_OPS_TIMEOUT: %v, set to %v", err, DEFAULT_ETCD_OPS_TIMEOUT)
	} else {
		kv.OprTimeout = opsTimeout
	}
	return &kv
}

func (kv *KvDepot) Get(key string) (*KV, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kv.OprTimeout)
	defer cancel()
	k := kv.Depot + "/" + key
	resp, err := kv.Conn.Get(ctx, k)
	if err != nil {
		e := fmt.Errorf("unable to get etcd key <<%s>> %v", k, err)
		kv.Log.Error(e)
		return nil, e
	}
	if resp.Kvs == nil || len(resp.Kvs) < 1 {
		kv.Log.Warnf("key %s not exist", k)
		return nil, nil
	}
	return &KV{
		Key:      string(resp.Kvs[0].Key),
		Value:    resp.Kvs[0].Value,
		Lease:    resp.Kvs[0].Lease,
		Revision: resp.Header.Revision,
	}, nil
}

func (kv *KvDepot) GetDir(path string) ([]*KV, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kv.OprTimeout)
	defer cancel()
	k := kv.Depot + "/" + path
	resp, err := kv.Conn.Get(ctx, k, etcd.WithPrefix())
	if err != nil {
		e := fmt.Errorf("unable to get etcd prefix <<%s>> %v", k, err)
		kv.Log.Error(e)
		return nil, e
	}
	if resp.Kvs == nil || len(resp.Kvs) < 1 {
		kv.Log.Warnf("prefix %s not exist", k)
		return nil, nil
	}
	dataSet := []*KV{}
	for _, d := range resp.Kvs {
		data := KV{
			Key:      string(d.Key),
			Value:    d.Value,
			Lease:    d.Lease,
			Revision: resp.Header.Revision,
		}
		dataSet = append(dataSet, &data)
	}
	return dataSet, nil
}

func (kv *KvDepot) Put(key, val string, lease int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), kv.OprTimeout)
	defer cancel()
	k := kv.Depot + "/" + key
	opt := []etcd.OpOption{}
	if lease > 4 {
		resp, err := kv.Conn.Grant(context.TODO(), lease)
		if err != nil {
			e := fmt.Errorf("unable to grant lease %v", err)
			kv.Log.Error(e)
			return e
		}
		opt = append(opt, etcd.WithLease(resp.ID))
	}
	_, err := kv.Conn.Put(ctx, k, val, opt...)
	if err != nil {
		e := fmt.Errorf("unable to save key %s %v", k, err)
		kv.Log.Error(e)
		return e
	}
	return nil
}

func (kv *KvDepot) Subscribe(key string) (*KV, etcd.WatchChan, error) {
	currentKV, err := kv.Get(key)
	if err != nil {
		return nil, nil, fmt.Errorf("unnable to get current value - %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	kv.Cancel = cancel
	return currentKV, kv.Conn.Watch(ctx, kv.Depot+"/"+key), nil
}

func (kv *KvDepot) KeepAlive(key string) error {

	return nil
}

func (kv *KvDepot) FindLatest(key string) (*etcd.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), kv.OprTimeout)
	defer cancel()
	k := kv.Depot + "/" + key
	rch := kv.Conn.Watch(ctx, k, etcd.WithPrefix(), etcd.WithRev(1)) //, etcd.WithSort(etcd.SortByModRevision, etcd.SortDescend)
	for wresp := range rch {
		for i := len(wresp.Events) - 1; i >= 0; i-- {
			if wresp.Events[i].Kv.Version != 0 {
				cancel()
				return wresp.Events[i], nil
			}
		}
		return nil, fmt.Errorf("invalid etcd watch event")
	}
	fmt.Println("++++++++++++ None +++++++++++++++")
	return nil, nil
}
