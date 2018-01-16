package main

import (
	"plugin"
	graphite "github.com/almariah/go-graphite-client"
	"github.com/golang/glog"
	"time"
	"fmt"
)

func RunPlugins(client graphite.Client, plugins []AgentPluginSpec, prefix string, hostname string) {

	for _, pluginSpec := range plugins {
		plug, err := plugin.Open("/opt/ik-agent/lib/" + pluginSpec.Name + ".so")
		if err != nil {
			glog.Fatal(err)
		}
		initSymbol, err := plug.Lookup("Init")
		if err != nil {
			glog.Fatal(err)
		}
		init := initSymbol.(func(map[string]interface{}))
		init(pluginSpec.Parameters)
		for _, metric  := range pluginSpec.Mertics {
			metricSymbol, err := plug.Lookup(metric.Name)
			if err != nil {
				glog.Fatal(err)
			}
			metricFunc := metricSymbol.(func(map[string]interface{}) string)
			go func(pluginName string, metricName string, res time.Duration, param map[string]interface{}){
				if prefix != "" {
					prefix = prefix + "."
				}
				for {
					m := graphite.Metric{
						Name:      prefix+hostname+"."+pluginName+"."+metricName,
						Value:     metricFunc(param),
						Timestamp: time.Now().Unix(),
					}
					err := client.SendMetric(m)
					if err != nil {
  					glog.Error(err)
  				}
					time.Sleep(res*time.Second)
					fmt.Println(m)
				}
			}(pluginSpec.Name, metric.Name, metric.Resolution, metric.Parameters)
		}
	}
}
