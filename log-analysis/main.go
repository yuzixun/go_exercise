package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"time"

	parser "log-analysis/parser"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type CmdParams struct {
	logFilePath string
	routineNum  int
}

type urlData struct {
	nginxLog parser.NgxLog
	uid      string
}

var logger = logrus.New()
var redisClient *redis.Pool

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.DebugLevel)
}

func main() {
	redisClient = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   10,
		IdleTimeout: 10 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", ":6379",
				redis.DialConnectTimeout(3*time.Second),
				redis.DialReadTimeout(3*time.Second),
				redis.DialWriteTimeout(3*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

	pLogFilePath := flag.String("logFilePath", "./access.log", "log file path")
	pRoutineNum := flag.Int("routineNum", 5, "consumer number of go routine")
	pLogPath := flag.String("logPath", "./process.log", "save log for this application")

	flag.Parse()

	params := CmdParams{logFilePath: *pLogFilePath, routineNum: *pRoutineNum}
	// fmt.Println("params is", params, *pLogPath)

	logFd, err := os.OpenFile(*pLogPath, os.O_CREATE|os.O_WRONLY, 0777)
	if err == nil {
		logger.Out = logFd
		defer logFd.Close()
	}

	logger.Infoln("application started")
	logger.Infoln("Params is ", params)

	logChannel := make(chan parser.NgxLog, 3*params.routineNum)
	pvChannel := make(chan urlData, params.routineNum)
	uvChannel := make(chan urlData, params.routineNum)

	go parser.Handle(params.logFilePath, logChannel)

	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	go pvCounter(pvChannel)
	go uvCounter(uvChannel)

	time.Sleep(time.Second * 1000)
}

func logConsumer(logChannel chan parser.NgxLog, pvChannel chan urlData, uvChannel chan urlData) {
	for ngxLog := range logChannel {
		// fmt.Println(ngxLog)

		hasher := md5.New()
		hasher.Write([]byte(ngxLog.RemoteAddr + ngxLog.HttpUserAgent))
		uid := hex.EncodeToString(hasher.Sum(nil))

		uData := urlData{ngxLog, uid}
		logger.Infoln(uData)
		pvChannel <- uData
		uvChannel <- uData
	}
}

func pvCounter(pvChannel chan urlData) {
	// 从池里获取连接
	redisCon := redisClient.Get()
	// 用完后将连接放回连接池
	defer redisCon.Close()
	for log := range pvChannel {

		// fmt.Println("pv", log.nginxLog.Time)
		time := log.nginxLog.Time.Format("2006010215")
		key := fmt.Sprintf("la:PV:%s", time)
		_, err := redisCon.Do("zincrby", key, 1, log.nginxLog.RequestURL)
		// fmt.Println("pv", key, log.nginxLog.RequestURL)
		if err != nil {
			logger.Errorln("pv set redis err", err)
		}
	}
}

func uvCounter(uvChannel chan urlData) {
	// 从池里获取连接
	redisCon := redisClient.Get()
	// 用完后将连接放回连接池
	defer redisCon.Close()
	for log := range uvChannel {

		// fmt.Println("uv", log.nginxLog.Time)
		time := log.nginxLog.Time.Format("20060102")
		hllKey := fmt.Sprintf("la:kll:%s", time)

		result, err := redisCon.Do("pfadd", hllKey, log.uid)
		// fmt.Println("uv", hllKey, result, log.uid)
		if err != nil {
			fmt.Println(err)
		}

		if result.(int64) != 1 {
			continue
		}

		key := fmt.Sprintf("la:UV:%s", time)
		_, err = redisCon.Do("zincrby", key, 1, log.nginxLog.RequestURL)
		// fmt.Println("uv", key, log.nginxLog.RequestURL)
		if err != nil {
			logger.Errorln("uv set redis err", err)
		}
	}
}
