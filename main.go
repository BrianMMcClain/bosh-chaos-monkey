package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func main() {

	// Parse CLI flags
	configPath := flag.String("config", "./config.json", "Path to config file")
	dryRun := flag.Bool("dry", false, "Dry Run")
	flag.Parse()

	conf := parseConfig(*configPath)

	// Read CA
	f, _ := os.Open(conf.CAPath)
	reader := bufio.NewReader(f)
	caContent, _ := ioutil.ReadAll(reader)

	// Connect to BOSH director
	director, err := buildDirector(conf.DirectorURL, conf.Username, conf.Password, string(caContent))
	if err != nil {
		panic(err)
	}

	// Enable resurrection
	rezerr := director.EnableResurrection(true)
	if rezerr != nil {
		panic(rezerr)
	}

	// Find deployment to test
	var monkeyDep boshdir.Deployment
	deps, err := director.Deployments()
	if err != nil {
		panic(err)
	}
	for _, dep := range deps {
		if dep.Name() == conf.DeploymentName {
			monkeyDep = dep
		}
	}

	if monkeyDep == nil {
		fmt.Printf("Could not find deployment named \"%s\", exiting\n", conf.DeploymentName)
		return
	}

	go monkey(monkeyDep, conf.KillInterval, *dryRun)

	// This probably isn't how you're supposed to do this
	for true {
		time.Sleep(time.Second * 15)
	}
}

func buildDirector(directorURL string, username string, password string, ca string) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(directorURL)
	if err != nil {
		return nil, err
	}

	config.Client = username
	config.ClientSecret = password
	config.CACert = ca

	return factory.New(config, boshdir.NewNoopTaskReporter(), boshdir.NewNoopFileReporter())
}

func monkey(deployment boshdir.Deployment, interval int, dryRun bool) {
	for true {
		time.Sleep(time.Second * time.Duration(interval))
		instances, _ := deployment.Instances()
		instance := instances[rand.Intn(len(instances))]
		fmt.Printf("Killing %s", instance.ID)
		if dryRun {
			fmt.Printf(" -- Skipping due to dry run\n")
		} else {
			deployment.DeleteVM(instance.VMID)
			fmt.Printf("\n")
		}
	}
}
