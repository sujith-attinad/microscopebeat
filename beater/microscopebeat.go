package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/sujith-attinad/microscopebeat/config"
	"os/exec"
	"strings"
	"strconv"
)

type Microscopebeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	config := config.DefaultConfig
	if err := cfg.Unpack(&config); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}
	bt := &Microscopebeat{
		done:   make(chan struct{}),
		config: config,
	}
	return bt, nil
}

func (bt *Microscopebeat) Run(b *beat.Beat) error {
	logp.Info("microscopebeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

        for _, table := range bt.config.Table {
        	bt.getLatency(table, counter, b)
			logp.Info("Event sent")
			counter++
		}
	}
}
func (bt *Microscopebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Microscopebeat) getLatency(table string, counter int, b *beat.Beat) {
	cmdCheck := fmt.Sprintf("nodetool cfstats %s | awk ' FNR == 2 { print $3 } FNR == 4 { print $3 } '", table)
	logp.Info(cmdCheck)
	outt, err := exec.Command("bash","-c",cmdCheck).Output()
	if err != nil {
		fmt.Sprintf("Failed to execute command: %s", cmdCheck)
	}
	logp.Info("hmmm" + string(outt))

	latency := strings.Split(string(outt), "\n")

	var read_latency, write_latency float64
	if strings.Compare(latency[0], "NaN") == 0 {
		read_latency = 0.0
	} else {
		read_latency, _ = strconv.ParseFloat(latency[0], 64)
	}
	if strings.Compare(latency[1], "NaN") == 0 {
		write_latency = 0.0
	} else {
		write_latency, _ = strconv.ParseFloat(latency[1], 64)
	}


	event := beat.Event{
		Timestamp: time.Now(),
		Fields: common.MapStr{
			"type":    b.Info.Name,
			"counter": counter,
			"table_name":	 table,
			"write_latency": write_latency,
			"read_latency":	 read_latency,
			},
	}


	bt.client.Publish(event)
}
