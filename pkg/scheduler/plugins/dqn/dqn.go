package dqn

import (
	"fmt"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"volcano.sh/volcano/pkg/scheduler/api"
	"volcano.sh/volcano/pkg/scheduler/framework"
)

const (
	// PluginName indicates name of volcano scheduler plugin.
	PluginName = "dqn"
)

/*
   actions: "enqueue, allocate"
   tiers:
   - plugins:
     - name: dqn
*/

type usagePlugin struct {
}

// New function returns usagePlugin object
func New(args framework.Arguments) framework.Plugin {
	return &usagePlugin{}
}

func (up *usagePlugin) Name() string {
	return PluginName
}

func (up *usagePlugin) OnSessionOpen(ssn *framework.Session) {
	klog.V(5).Infof("Enter dqn plugin ...")
	defer func() {
		klog.V(5).Infof("Leaving dqn plugin ...")
	}()

	if klog.V(4).Enabled() {
		for node, nodeInfo := range ssn.Nodes {
			klog.V(4).Infof("node:%v, cpu usage:%v, mem usage:%v, metrics time is %v",
				node, nodeInfo.ResourceUsage.CPUUsageAvg, nodeInfo.ResourceUsage.MEMUsageAvg, nodeInfo.ResourceUsage.MetricsTime)
		}
	}
	bestNodeFn := func(task *api.TaskInfo, nodeScores map[float64][]*api.NodeInfo) *api.NodeInfo {
		choose := ""
		schedulerUrl := "http://192.168.3.108:1234/choose"
		urlValues := url.Values{}
		urlValues.Add("podname", task.Pod.Name)
		resp, err := http.PostForm(schedulerUrl, urlValues)
		if err != nil {
			fmt.Printf("[ERROR] Get choose from %v failed!\n\n", schedulerUrl)
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			choose = string(body)
			fmt.Printf("[INFO] Get choose from DQN: %v\n\n", choose)
			for _, nodes := range nodeScores {
				for _, node := range nodes {
					if node.Node.Name == choose {
						return node
					}
				}
			}
		}
		return nil
	}
	ssn.AddBestNodeFn(up.Name(), bestNodeFn)
}

func (up *usagePlugin) OnSessionClose(ssn *framework.Session) {}
