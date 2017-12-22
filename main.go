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
	directorURL := flag.String("director", "", "BOSH Directory URL")
	username := flag.String("username", "", "BOSH Username")
	password := flag.String("password", "", "BOSH Password")
	caPath := flag.String("ca", "", "Path to CA Cert")
	deploymentName := flag.String("deployment", "", "Name of deployment")
	interval := flag.Int("interval", 60, "Chaos interval to kill machines, in seconds")
	dryRun := flag.Bool("dry", false, "Dry Run")
	flag.Parse()

	// Connect to BOSH director
	director, err := buildDirector(*directorURL, *username, *password, *caPath)
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
		if dep.Name() == *deploymentName {
			monkeyDep = dep
		}
	}

	if monkeyDep == nil {
		fmt.Printf("Could not find deployment named \"%s\", exiting\n", *deploymentName)
		return
	}

	go monkey(monkeyDep, *interval, *dryRun)

	// This probably isn't how you're supposed to do this
	for true {
		time.Sleep(time.Second * 15)
	}
}

func buildDirector(directorURL string, username string, password string, caPath string) (boshdir.Director, error) {
	logger := boshlog.NewLogger(boshlog.LevelError)
	factory := boshdir.NewFactory(logger)

	config, err := boshdir.NewConfigFromURL(directorURL)
	if err != nil {
		return nil, err
	}

	// Read CA
	f, _ := os.Open(caPath)
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	config.Client = username
	config.ClientSecret = password
	config.CACert = string(content)

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
