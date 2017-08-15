package beater

import (
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher"

	"github.com/gamegos/etcdbeat/config"
)

type Etcdbeat struct {
	done         chan struct{}
	config       config.EtcdbeatConfig
	client       publisher.Client
	auth         bool
	authEnable   bool
	username     string
	password     string
	EbConfig     config.ConfigSettings
	period       time.Duration
	port         string
	host         string
	leaderEnable bool
	selfEnable   bool
	storeEnable  bool
}

const selector = "etcdbeat"

func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	eb := &Etcdbeat{
		done: make(chan struct{}),
	}
	err := cfgfile.Read(&eb.EbConfig, "")
	if err != nil {
		logp.Err("Error reading configuration file: %v", err)
		return nil, fmt.Errorf("Error reading configuration file: %v", err)
	}
	return eb, nil
}

func (eb *Etcdbeat) StatisticsCheck(b *beat.Beat) {

	if eb.leaderEnable {
		leaderstats, err := eb.getLeaderStats(b)
		if err != nil {
			logp.Debug(selector, "Error reading leader stats")
		} else {
			eventleader := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"leader":     leaderstats,
			}
			eb.client.PublishEvent(eventleader)
			logp.Info("Leader Stats: event sent")
		}
	}

	if eb.selfEnable {
		selfstats, err := eb.getSelfStats(b)
		if err != nil {
			logp.Debug(selector, "Error self leader stats")
		} else {
			eventself := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"self":       selfstats,
			}
			eb.client.PublishEvent(eventself)
			logp.Info("Self Stats: Event sent")
		}
	}

	if eb.storeEnable {
		storestats, err := eb.getStoreStats(b)
		if err != nil {
			logp.Debug(selector, "Error reading store stats")
		} else {
			eventstore := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Name,
				"store":      storestats,
			}
			eb.client.PublishEvent(eventstore)
			logp.Info("Store Stats: Event sent")
		}
	}
}

func (eb *Etcdbeat) Run(b *beat.Beat) error {

	logp.Info("etcdbeat is running! Hit CTRL-C to stop it.")
	eb.CheckConfig(b)

	eb.client = b.Publisher.Connect()
	ticker := time.NewTicker(eb.period)
	defer ticker.Stop()
	for {
		select {
		case <-eb.done:
			return nil
		case <-ticker.C:
		}

		if eb.authEnable {
			if eb.auth {
				eb.StatisticsCheck(b)
			} else {
				logp.Debug(selector, "Username or Password not set.")
			}
		} else {
			eb.StatisticsCheck(b)
		}
	}
}

func (eb *Etcdbeat) CheckConfig(b *beat.Beat) error {

	if eb.EbConfig.Input.Period != nil {
		eb.period = time.Duration(*eb.EbConfig.Input.Period) * time.Second
	} else {
		eb.period = 30 * time.Second
	}

	if eb.EbConfig.Input.Port != nil {
		eb.port = *eb.EbConfig.Input.Port
	} else {
		eb.port = "2379"
	}

	if eb.EbConfig.Input.Host != nil {
		eb.host = *eb.EbConfig.Input.Host
	} else {
		eb.port = "localhost"
	}

	eb.authEnable = *eb.EbConfig.Input.Authentication.Enable

	if eb.authEnable {
		if eb.EbConfig.Input.Authentication.Username == nil || eb.EbConfig.Input.Authentication.Password == nil {
			logp.Err("Username or password is not set.")
			eb.auth = false
		} else if *eb.EbConfig.Input.Authentication.Username == "" || *eb.EbConfig.Input.Authentication.Password == "" {
			logp.Err("Username or password is not set.")
			eb.auth = false
		} else {
			eb.username = *eb.EbConfig.Input.Authentication.Username
			eb.password = *eb.EbConfig.Input.Authentication.Password
			eb.auth = true
			logp.Debug(selector, "Username %v", eb.username)
			logp.Debug(selector, "Password %v", eb.password)

			req, err := http.NewRequest("GET", "http://"+eb.host+eb.port+"/v2/keys", nil)
			if err != nil {
				logp.Err("Error connect to port: %v", err)
			}
			req.SetBasicAuth(eb.username, eb.password)
			cli := &http.Client{}
			resp, err := cli.Do(req)
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusUnauthorized {
				eb.auth = false
				logp.Err("Username or password is wrong.")
			} else {
				eb.auth = true

			}
		}
	}

	eb.leaderEnable = *eb.EbConfig.Input.Statistics.Leader
	eb.selfEnable = *eb.EbConfig.Input.Statistics.Self
	eb.storeEnable = *eb.EbConfig.Input.Statistics.Store

	logp.Debug(selector, "Leader Statistic Enable = %v", eb.leaderEnable)
	logp.Debug(selector, "Self Statistic Enable = %v", eb.selfEnable)
	logp.Debug(selector, "Store Statistic Enable = %v", eb.storeEnable)

	logp.Debug(selector, "Init Etcdbeat")
	logp.Debug(selector, "Port %v", eb.port)
	logp.Debug(selector, "Period %v", eb.period)
	logp.Debug(selector, "Host %v", eb.host)

	return nil

}

func (eb *Etcdbeat) Stop() {
	eb.client.Close()
	close(eb.done)
}
