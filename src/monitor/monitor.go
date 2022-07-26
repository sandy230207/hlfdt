package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
)

func status(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("src/monitor/index.html"))
	data := new(IndexData)
	data.Title = "Status of Hyperledger Fabric Network"
	res := ""
	cmd := exec.Command("docker", "stats", "--no-stream")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Server error: ", err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			res += scanner.Text() + "\n"
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal("Server error: ", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal("Server error: ", err)
	}
	fmt.Println(res)
	data.Content = res

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	body := buf.String()
	body = strings.Replace(body, "\n", "<br>", -1)
	fmt.Fprint(w, body)
}

func blockHeight(w http.ResponseWriter, r *http.Request) {
	// peer channel fetch newest mychannel.block  -c mychannel
	tmpl := template.Must(template.ParseFiles("src/monitor/index.html"))
	data := new(IndexData)
	data.Title = "Block Height of Hyperledger Fabric Network"
	res := ""
	cmd := exec.Command("./network.sh", "monitor")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Server error: ", err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			res += scanner.Text() + "\n"
		}
	}()
	if err := cmd.Start(); err != nil {
		log.Fatal("Server error: ", err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal("Server error: ", err)
	}
	fmt.Println(res)
	data.Content = res

	buf := &bytes.Buffer{}
	if err := tmpl.Execute(buf, data); err != nil {
		panic(err)
	}
	body := buf.String()
	body = strings.Replace(body, "\n", "<br>", -1)
	fmt.Fprint(w, body)
}

func RunMonitor(c *cli.Context) {
	port := c.Args().Get(1)
	http.HandleFunc("/", status)
	http.HandleFunc("/index", status)
	http.HandleFunc("/block", blockHeight)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
