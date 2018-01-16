package main

import (
	"fmt"
	"os"
	//"time"
	graphite "github.com/almariah/go-graphite-client"
	//"strconv"
	//"github.com/almariah/ik-agent/core"
	"flag"
	"log"
	"github.com/ghodss/yaml"
	"io/ioutil"
	"time"
	"github.com/golang/glog"
)

type AgentConfig struct {
	InitConfig InitConfigSpec    `json:"initConfig"`
	Plugins    []AgentPluginSpec `json:"plugins"`
	Scripts    []string          `json:"scripts"`
}

type AgentPluginSpec struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
	Mertics    []MerticSpec           `json:"mertics"`
}

type MerticSpec struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
	Resolution time.Duration          `json:"resolution"`
}

type InitConfigSpec struct {
	GraphiteHost string `json:"GraphiteHost"`
	GraphitePort int    `json:"GraphitePort"`
	Prefix       string	`json:"prefix"`
}

func main() {

	confPtr := flag.String("c", "/opt/ik-agent/etc/ik-agent.yaml", "InfraKeeper agent config file location. (Required)")
	flag.Parse()

	glog.Info("Starting InfraKeeper Agent...")

	yamlConfig, err := ioutil.ReadFile(*confPtr)

	if err != nil {
		log.Fatal(err)
	}

	var config AgentConfig

	err = yaml.Unmarshal(yamlConfig, &config)
	if err != nil {
		log.Fatal("unable to parse config file: %v\n", err)
	}

	gClient := &graphite.Client{
		Host:     config.InitConfig.GraphiteHost,
		Port:     config.InitConfig.GraphitePort,
		Protocol: "tcp",
	}

	glog.Info("Connecting to graphite...")
	err = gClient.Connect()
	if err != nil {
		log.Fatal(err)
	}

	hostname, _ := os.Hostname()

	if config.Plugins != nil {
		RunPlugins(*gClient, config.Plugins, config.InitConfig.Prefix, hostname)
	}

	for {

	}

}

// exitWithError will terminate execution with an error result
// It prints the error to stderr and exits with a non-zero exit code
func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "\n%v\n", err)
	os.Exit(1)
}


/*
Config(*client)
for {
	select {
	case <-time.After(5*time.Second):
			{
				cpu := strconv.FormatFloat(core.Get(), 'f', 6, 64)
				mem := strconv.FormatInt(core.MemAvail(), 10)
				fmt.Println(cpu +":"+ mem)
				prefix := ""
				hostname, _ := os.Hostname()
				client.Send(prefix + hostname + ".cpu.usage", cpu)
				client.Send(prefix + hostname + ".memory.free", mem)
			}
	}
}*/
