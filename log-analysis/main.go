package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type CmdParams struct {
	logFilePath string
	routineNum  int
}

type LogData struct {
	remoteAddr           string
	remoteUser           string
	timeLocal            string
	request              string
	host                 string
	status               string
	requestLength        string
	bodyBytesSent        string
	httpReferer          string
	httpUserAgent        string
	forwardedFor         string
	upstreamAddr         string
	requestTime          string
	upstreamResponseTime string
}

type urlData struct {
	data    LogData
	uid     string
	urlNode urlNode
}

type urlNode struct {
	nodeType      string
	nodeRequestID int
	nodeURL       string
	nodeTime      string
}

type storageBlock struct {
	counterType  string // PV or Uv
	storageModel string
	unode        urlNode
}

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	//
	pLogFilePath := flag.String("logFilePath", "./access.log", "log file path")
	pRoutineNum := flag.Int("routineNum", 5, "consumer number of go routine")
	pLogPath := flag.String("logPath", "./logs", "save log for this application")

	flag.Parse()

	params := CmdParams{logFilePath: *pLogFilePath, routineNum: *pRoutineNum}
	fmt.Println("params ios", params, *pLogPath)
	//
	logFd, err := os.OpenFile(*pLogPath, os.O_CREATE|os.O_WRONLY, 0777)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	fmt.Println("xxxxxxxxxx")
	log.Infoln("application started")
	log.Infoln("Params is ", params)
	logChannel := make(chan string, 3*params.routineNum)
	pvChannel := make(chan urlData, params.routineNum)
	uvChannel := make(chan urlData, params.routineNum)
	storageChannel := make(chan storageBlock, params.routineNum)

	go parseLog(params, logChannel)

	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel)

	go dataPersist(storageChannel)

	// time.Sleep(time.Second * 1000)
	time.Sleep(time.Second * 3)
}

func dataPersist(storageChannel chan storageBlock) {

}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {
}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {

}

func logConsumer(logChannel chan string, pvChannel chan urlData, uvChannel chan urlData) {
	for line := range logChannel {
		// fmt.Println(line)
		data := parseLine(line)
		// fmt.Println(data)

		hasher := md5.New()
		hasher.Write([]byte(data.httpReferer + data.httpUserAgent))
		uid := hex.EncodeToString(hasher.Sum(nil))

		uData := urlData{*data, uid, formatURL((*data).request, (*data).timeLocal)}
		fmt.Println(uData)
		log.Infoln(uData)
		pvChannel <- uData
		uvChannel <- uData
	}
	// line := <-logChannel

}

func formatURL(request string, time string) urlNode {
	return urlNode{
		// request: "movie",
		nodeType:      "list",
		nodeRequestID: 0,
		nodeURL:       "url",
		nodeTime:      "time",
	}
}

func parseLog(params CmdParams, logChannel chan string) error {
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		log.Errorf("file can not open")
		return err
	}

	defer fd.Close()
	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(3 * time.Second)
				log.Infof("wait for three seconds")
			} else {
				log.Errorf("read Line Error")
			}
		}

		logChannel <- line
		break
		count++
		if count%(1000*params.routineNum) == 0 {
			log.Infof("read file line: %d", count)
		}
	}
	return nil
}
