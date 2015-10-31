package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/davecheney/profile"
	docker "github.com/fsouza/go-dockerclient"
)

var (
	before map[string]int64 = map[string]int64{}
	after  map[string]int64 = map[string]int64{}
	times  map[string]int64 = map[string]int64{}
)

const (
	noLabelContainerNumber = 5000
	labelContainerNumber   = 100
	labelKey               = "test"
	labelValue             = "1"
	singleLabelKey         = "single"
	singleLabelValue       = "1"
	listTestTimes          = 100
)

func createTestContainers(client *docker.Client) {
	before["create_container"] = time.Now().UnixNano()
	createNoLabelContainers(client, noLabelContainerNumber)
	createLabelContainers(client, labelContainerNumber)
	createSingleLabelContaienr(client)
	after["create_container"] = time.Now().UnixNano()
	times["create_container"] = noLabelContainerNumber + labelContainerNumber + 1
	printTimeCost("create_container")
}

func createNoLabelContainers(client *docker.Client, num int) {
	createContainers(client, num, map[string]string{})
}

func createLabelContainers(client *docker.Client, num int) {
	createContainers(client, num, map[string]string{labelKey: labelValue})
}

func createSingleLabelContaienr(client *docker.Client) {
	createContainers(client, 1, map[string]string{singleLabelKey: singleLabelValue})
}

func createContainers(client *docker.Client, num int, labels map[string]string) {
	for i := 0; i < num; i++ {
		name := "test_container_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		dockerOpts := docker.CreateContainerOptions{
			Name: name,
			Config: &docker.Config{
				Image:  "ubuntu",
				Labels: labels,
			},
		}
		client.CreateContainer(dockerOpts)
	}
}

func printTimeCost(item string) {
	fmt.Printf("%v time cost: %v(us)\n", item, (after[item]-before[item])/1000/times[item])
}

func main() {
	defer profile.Start(profile.CPUProfile).Stop()

	before["whole"] = time.Now().UnixNano()
	endpoint := "unix:///var/run/docker.sock"
	before["create_client"] = time.Now().UnixNano()
	client, _ := docker.NewClient(endpoint)
	after["create_client"] = time.Now().UnixNano()
	times["create_client"] = 1
	printTimeCost("create_client")

	// createTestContainers(client)
	before["list_containers"] = time.Now().UnixNano()
	for i := 0; i < listTestTimes; i++ {
		client.ListContainers(docker.ListContainersOptions{All: true})
	}
	after["list_containers"] = time.Now().UnixNano()
	times["list_containers"] = listTestTimes
	printTimeCost("list_containers")

	before["filter_containers"] = time.Now().UnixNano()
	for i := 0; i < listTestTimes; i++ {
		client.ListContainers(docker.ListContainersOptions{All: true, Filters: map[string][]string{"label": []string{labelKey + "=" + labelValue}}})
	}
	after["filter_containers"] = time.Now().UnixNano()
	times["filter_containers"] = listTestTimes
	printTimeCost("filter_containers")

	before["single_container"] = time.Now().UnixNano()
	for i := 0; i < listTestTimes; i++ {
		client.ListContainers(docker.ListContainersOptions{All: true, Filters: map[string][]string{"label": []string{singleLabelKey + "=" + singleLabelValue}}})
	}
	after["single_container"] = time.Now().UnixNano()
	times["single_container"] = listTestTimes
	printTimeCost("single_container")

	after["whole"] = time.Now().UnixNano()
	times["whole"] = 1
	printTimeCost("whole")
}
