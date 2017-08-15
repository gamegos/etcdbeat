package beater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
)

type LeaderStats struct {
	Followers map[string]struct {
		Counts struct {
			Fail    int `json:"fail"`
			Success int `json:"success"`
		} `json:"counts"`
		Latency struct {
			Average           float64 `json:"average"`
			Current           float64 `json:"current"`
			Maximum           float64 `json:"maximum"`
			Minimum           int     `json:"minimum"`
			StandardDeviation float64 `json:"standardDeviation"`
		} `json:"latency"`
	} `json:"followers"`
	Leader string `json:"leader"`
}

type SelfStats struct {
	ID         string `json:"id"`
	LeaderInfo struct {
		Leader    string `json:"leader"`
		StartTime string `json:"startTime"`
		Uptime    string `json:"uptime"`
	} `json:"leaderInfo"`
	Name                 string  `json:"name"`
	RecvAppendRequestCnt int     `json:"recvAppendRequestCnt"`
	RecvBandwidthRate    float64 `json:"recvBandwidthRate"`
	RecvPkgRate          float64 `json:"recvPkgRate"`
	SendAppendRequestCnt int     `json:"sendAppendRequestCnt"`
	StartTime            string  `json:"startTime"`
	State                string  `json:"state"`
}

type StoreStats struct {
	GetsSuccess             int `json:"getsSuccess"`
	GetsFail                int `json:"getsFail"`
	SetsSuccess             int `json:"setsSuccess"`
	SetsFail                int `json:"setsFail"`
	DeleteSuccess           int `json:"deleteSuccess"`
	DeleteFail              int `json:"deleteFail"`
	UpdateSuccess           int `json:"updateSuccess"`
	UpdateFail              int `json:"updateFail"`
	CreateSuccess           int `json:"createSuccess"`
	CreateFail              int `json:"createFail"`
	CompareAndSwapSuccess   int `json:"compareAndSwapSuccess"`
	CompareAndSwapFail      int `json:"compareAndSwapFail"`
	CompareAndDeleteSuccess int `json:"compareAndDeleteSuccess"`
	CompareAndDeleteFail    int `json:"compareAndDeleteFail"`
	ExpireCount             int `json:"expireCount"`
	Watchers                int `json:"watchers"`
}

func (eb *Etcdbeat) connectAPI(url string, statsname string) ([]uint8, error) {
	res, err := http.Get("http://" + eb.host + eb.port + url)
	if err != nil {
		logp.Err("%q = Error connecting Etcd: %v", statsname, err)
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		logp.Err("API returned wrong status code: HTTP %s ", res.Status)
		return nil, fmt.Errorf("HTTP %s", res.Status)
	}

	resp, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		logp.Err("Error reading stats: %v", err)
		return nil, fmt.Errorf("HTTP%s", res.Status)
	}
	return resp, nil
}

func (eb *Etcdbeat) getLeaderStats(b *beat.Beat) (LeaderStats, error) {
	url := "/v2/stats/leader"
	var response LeaderStats

	resp, err := eb.connectAPI(url, "leader")
	if err != nil {
		return response, err
	}
	unmarshalErr := json.Unmarshal(resp, &response)
	if unmarshalErr != nil {
		logp.Err(unmarshalErr.Error())
		return response, unmarshalErr
	}
	return response, nil
}

func (eb *Etcdbeat) getSelfStats(b *beat.Beat) (SelfStats, error) {
	url := "/v2/stats/self"
	var response SelfStats

	resp, err := eb.connectAPI(url, "self")
	if err != nil {
		return response, err
	}

	unmarshalErr := json.Unmarshal(resp, &response)
	if unmarshalErr != nil {
		logp.Err(unmarshalErr.Error())
		return response, unmarshalErr
	}
	return response, nil

}

func (eb *Etcdbeat) getStoreStats(b *beat.Beat) (StoreStats, error) {
	url := "/v2/stats/store"
	var response StoreStats

	resp, err := eb.connectAPI(url, "store")
	if err != nil {
		return response, err
	}
	unmarshalErr := json.Unmarshal(resp, &response)
	if unmarshalErr != nil {
		logp.Err(unmarshalErr.Error())
		return response, unmarshalErr
	}
	return response, nil
}
