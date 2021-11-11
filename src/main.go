package main

import (
	"bufio"
	"fabric-tool/src/config"
	"fabric-tool/src/configtx"
	dockerCA "fabric-tool/src/docker-ca"
	dockerCouch "fabric-tool/src/docker-couch"
	dockerNet "fabric-tool/src/docker-net"
	serverconfig "fabric-tool/src/fabric-ca-server-config"
	fabricNetwork "fabric-tool/src/fabric-network"
	"fabric-tool/src/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func createNetwork(conf *config.Config) error {
	for _, channel := range conf.Channels {
		cmd := exec.Command("./network.sh", "up", "createChannel", "-c", channel.Name, "-ca")
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Println(scanner.Text())
			}
		}()
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
	}

	for _, cc := range conf.Chaincodes {
		for _, channel := range cc.Channels {
			cmd := exec.Command("./network.sh", "deployCC", "-c", channel, "-ccn", cc.Name, "-ccp", cc.Path, "-ccl", cc.Language)
			cmdReader, err := cmd.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			scanner := bufio.NewScanner(cmdReader)
			go func() {
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
			}()
			if err := cmd.Start(); err != nil {
				log.Fatal(err)
			}
			if err := cmd.Wait(); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func main() {
	file := os.Args[1]

	fmt.Println("Parsing Config File...")
	rConf, err := config.ReadConf(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Checking Config...")
	rConf, err = config.CheckConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	txConf, err := configtx.ConvertConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating configtx.yaml...")
	err = utils.WriteYaml(txConf, "./configtx/configtx.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = utils.ConvertConfigtx("./configtx/configtx.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	caConf, err := dockerCA.ConvertConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating fabric-ca-server-config.yaml for all organizations...")
	err = serverconfig.MakeDirsAndWriteConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating docker-compose-ca.yaml...")
	err = utils.WriteYaml(caConf, "./docker/docker-compose-ca.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = utils.ConvertNet("./docker/docker-compose-ca.yaml", "networks:", "services:")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	couchConf, err := dockerCouch.ConvertConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating docker-compose-couch.yaml...")
	err = utils.WriteYaml(couchConf, "./docker/docker-compose-couch.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = utils.ConvertNet("./docker/docker-compose-couch.yaml", "networks:", "services:")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	netConf, err := dockerNet.ConvertConf(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating docker-compose-net.yaml...")
	err = utils.WriteYaml(netConf, "./docker/docker-compose-test-net.yaml")
	if err != nil {
		log.Fatal(err)
	}
	err = utils.ConvertNet("./docker/docker-compose-test-net.yaml", "volumes:", "networks:")
	if err != nil {
		log.Fatal(err)
	}
	err = utils.ConvertNet("./docker/docker-compose-test-net.yaml", "networks:", "services:")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	reConf, err := fabricNetwork.GenerateEnrollRegister(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating registerEnroll.sh...")
	err = utils.WriteSh(reConf, "./organizations/fabric-ca/registerEnroll.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	networkConf, err := fabricNetwork.GenerateNetwork(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating network.sh...")
	err = utils.WriteSh(networkConf, "./network.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	deployCCConf, err := fabricNetwork.GenerateDeployCC(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating deployCC.sh...")
	err = utils.WriteSh(deployCCConf, "./scripts/deployCC.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	envVarConf, err := fabricNetwork.GenerateEnvVar(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating envVar.sh...")
	err = utils.WriteSh(envVarConf, "./scripts/envVar.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	configUpdateConf, err := fabricNetwork.GenerateConfigUpdate(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating configUpdate.sh...")
	err = utils.WriteSh(configUpdateConf, "./scripts/configUpdate.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	configSerAnchorPeer, err := fabricNetwork.GenerateSetAnchorPeer(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating setAnchorPeer.sh...")
	err = utils.WriteSh(configSerAnchorPeer, "./scripts/setAnchorPeer.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	configCreateChannel, err := fabricNetwork.GenerateCreateChannel(rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating createChannel.sh...")
	err = utils.WriteSh(configCreateChannel, "./scripts/createChannel.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Analysing Config...")
	configCCPGenerate, err := fabricNetwork.GenerateCCPGenerate("./organizations/ccp-generate-script.sh", rConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Generating ccp-generate.sh...")
	err = utils.WriteSh(configCCPGenerate, "./organizations/ccp-generate.sh")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Scripts written.")
	fmt.Println("Build & Deploy Hyperledger Fabric Network...")
	// createNetwork(rConf)

	// fmt.Println(rConf.Organizations[0].Peers)
	// fmt.Println(caConf)
}
