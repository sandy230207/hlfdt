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
	"fabric-tool/src/monitor"
	"fabric-tool/src/utils"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name: "hlfdt",
		Commands: []cli.Command{
			{
				Name:            "generate",
				Aliases:         []string{"g"},
				Usage:           "Generate the config files and scripts",
				Action:          generate,
				SkipFlagParsing: true,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "file",
						Value: "./config.yaml",
						Usage: "The directory of target config file used for generating Hyperledger Fabric Network",
					},
				},
			},
			{
				Name:   "up",
				Usage:  "Create and start Hyperledger Fabric Network",
				Action: networkUp,
			},
			{
				Name:   "down",
				Usage:  "Stop and clear Hyperledger Fabric Network",
				Action: networkDown,
			},
			{
				Name:            "ui",
				Usage:           "Lookup the status of Hyperledger Fabric Network",
				Action:          monitor.RunMonitor,
				SkipFlagParsing: true,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "port",
						Value: "8888",
						Usage: "The port number that monitor will use",
					},
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func networkUp(c *cli.Context) error {
	cmd := exec.Command("./network.sh", "up", "createChannel", "-ca")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	if err := deployCC(); err != nil {
		return err
	}
	return nil
}

func deployCC() error {
	cmd := exec.Command("./network.sh", "deployCC")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func networkDown(c *cli.Context) error {
	cmd := exec.Command("./network.sh", "down")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func generate(c *cli.Context) error {
	file := c.Args().Get(1)

	fmt.Println("Parsing Config File...")
	rConf, err := config.ReadConf(file)
	if err != nil {
		return err
	}
	fmt.Println("Checking Config...")
	rConf, err = config.CheckConf(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	for _, channel := range rConf.Channels {
		txConf, err := configtx.ConvertConf(rConf, channel.Name)
		if err != nil {
			return err
		}
		fmt.Println("Generating configtx.yaml...")
		_ = os.Mkdir("./configtx-"+channel.Name, 0777)
		err = utils.WriteYaml(txConf, "./configtx-"+channel.Name+"/configtx.yaml")
		if err != nil {
			return err
		}
	}

	err = utils.ConvertConfigtx("./configtx/configtx.yaml")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	caConf, err := dockerCA.ConvertConf(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating fabric-ca-server-config.yaml for all organizations...")
	err = serverconfig.MakeDirsAndWriteConf(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating docker-compose-ca.yaml...")
	err = utils.WriteYaml(caConf, "./docker/docker-compose-ca.yaml")
	if err != nil {
		return err
	}
	err = utils.ConvertNet("./docker/docker-compose-ca.yaml", "networks:", "services:")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	couchConf, err := dockerCouch.ConvertConf(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating docker-compose-couch.yaml...")
	err = utils.WriteYaml(couchConf, "./docker/docker-compose-couch.yaml")
	if err != nil {
		return err
	}
	err = utils.ConvertNet("./docker/docker-compose-couch.yaml", "networks:", "services:")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	netConf, err := dockerNet.ConvertConf(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating docker-compose-net.yaml...")
	err = utils.WriteYaml(netConf, "./docker/docker-compose-test-net.yaml")
	if err != nil {
		return err
	}
	err = utils.ConvertNet("./docker/docker-compose-test-net.yaml", "volumes:", "networks:")
	if err != nil {
		return err
	}
	err = utils.ConvertNet("./docker/docker-compose-test-net.yaml", "networks:", "services:")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	reConf, err := fabricNetwork.GenerateEnrollRegister(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating registerEnroll.sh...")
	err = utils.WriteSh(reConf, "./organizations/fabric-ca/registerEnroll.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	networkConf, err := fabricNetwork.GenerateNetwork(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating network.sh...")
	err = utils.WriteSh(networkConf, "./network.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	deployCCConf, err := fabricNetwork.GenerateDeployCC(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating deployCC.sh...")
	err = utils.WriteSh(deployCCConf, "./scripts/deployCC.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	envVarConf, err := fabricNetwork.GenerateEnvVar(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating envVar.sh...")
	err = utils.WriteSh(envVarConf, "./scripts/envVar.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	configUpdateConf, err := fabricNetwork.GenerateConfigUpdate(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating configUpdate.sh...")
	err = utils.WriteSh(configUpdateConf, "./scripts/configUpdate.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	configSerAnchorPeer, err := fabricNetwork.GenerateSetAnchorPeer(rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating setAnchorPeer.sh...")
	err = utils.WriteSh(configSerAnchorPeer, "./scripts/setAnchorPeer.sh")
	if err != nil {
		return err
	}
	fmt.Println("Analysing Config...")
	for _, channel := range rConf.Channels {
		configCreateChannel, err := fabricNetwork.GenerateCreateChannel(rConf, channel.Name)
		if err != nil {
			return err
		}
		fmt.Println("Generating createChannel.sh...")
		err = utils.WriteSh(configCreateChannel, "./scripts/createChannel-"+channel.Name+".sh")
		if err != nil {
			return err
		}
	}

	fmt.Println("Analysing Config...")
	configCCPGenerate, err := fabricNetwork.GenerateCCPGenerate("./organizations/ccp-generate-script.sh", rConf)
	if err != nil {
		return err
	}
	fmt.Println("Generating ccp-generate.sh...")
	err = utils.WriteSh(configCCPGenerate, "./organizations/ccp-generate.sh")
	if err != nil {
		return err
	}
	fmt.Println("Scripts written.")
	return nil
}
