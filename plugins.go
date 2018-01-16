package main

import (
	"plugin"
	graphite "github.com/almariah/go-graphite-client"
	"github.com/golang/glog"
	"time"
	//"fmt"
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

		init := initSymbol.(func(map[string]interface{}) error)

		err = init(pluginSpec.Parameters)
		if err != nil {
			glog.Fatal(err)
		}

		for _, metric  := range pluginSpec.Mertics {

			metricSymbol, err := plug.Lookup(metric.Name)
			if err != nil {
				glog.Fatal(err)
			}

			metricFunc := metricSymbol.(func(map[string]interface{}) (string, error))

			go func(pluginName string, metricName string, res time.Duration, param map[string]interface{}){
				if prefix != "" {
					prefix = prefix + "."
				}
				for {
					value, err := metricFunc(param)
					if err != nil {
						glog.Error(err)
						time.Sleep(res*time.Second)
						continue
					}
					m := graphite.Metric{
						Name:      prefix+hostname+"."+pluginName+"."+metricName,
						Value:     value,
						Timestamp: time.Now().Unix(),
					}
					err = client.SendMetric(m)
					if err != nil {
  					glog.Error(err)
  				}
					time.Sleep(res*time.Second)
				}
			}(pluginSpec.Name, metric.Name, metric.Resolution, metric.Parameters)

		}
	}
}
