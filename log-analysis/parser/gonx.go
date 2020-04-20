package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/satyrius/gonx"
)

type NgxLog struct {
	remoteAddr           string
	remoteUser           string
	timeLocal            string
	request              string
	httpHost             string
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

const NGX_FORMAT = "$remote_addr  $remote_user [$time_local] \"$request\" \"$http_host\" $status $request_length $body_bytes_sent \"$http_referer\" \"$http_user_agent\"     $http_x_forwarded_for   \"$upstream_addr\" \"$request_time\"        \"$upstream_response_time\""

func Handle(logFilePath string, logProcessor chan NgxLog) {
	logReader, err := os.Open(logFilePath)

	if err != nil {
		panic(err)
	}

	defer logReader.Close()

	reader := gonx.NewReader(logReader, NGX_FORMAT)

	for {
		rec, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				break
			}

			panic(err)
		}
		fmt.Printf("%+v\n", rec)

		logProcessor <- generateNgxLogBy(rec)
	}
}

func generateNgxLogBy(record *gonx.Entry) NgxLog {
	var err error
	log := NgxLog{}

	log.bodyBytesSent, err = record.Field("body_bytes_sent")
	if err != nil {
		log.bodyBytesSent = ""
		fmt.Println("parse bodyBytesSent failed.")
	}

	log.httpHost, err = record.Field("http_host")
	if err != nil {
		log.httpHost = ""
		fmt.Println("parse httpHost failed.")
	}

	log.httpReferer, err = record.Field("http_referer")
	if err != nil {
		log.httpReferer = ""
		fmt.Println("parse httpReferer failed.")
	}

	log.httpUserAgent, err = record.Field("http_user_agent")
	if err != nil {
		log.httpUserAgent = ""
		fmt.Println("parse httpUserAgent failed.")
	}

	log.forwardedFor, err = record.Field("http_x_forwarded_for")
	if err != nil {
		log.forwardedFor = ""
		fmt.Println("parse forwardedFor failed.")
	}

	log.remoteAddr, err = record.Field("remote_addr")
	if err != nil {
		log.remoteAddr = ""
		fmt.Println("parse remoteAddr failed.")
	}

	log.remoteUser, err = record.Field("remote_user")
	if err != nil {
		log.remoteUser = ""
		fmt.Println("parse remoteUser failed.")
	}

	log.request, err = record.Field("request")
	if err != nil {
		log.request = ""
		fmt.Println("parse request failed.")
	}

	log.requestLength, err = record.Field("request_length")
	if err != nil {
		log.requestLength = ""
		fmt.Println("parse requestLength failed.")
	}

	log.requestTime, err = record.Field("request_time")
	if err != nil {
		log.requestTime = ""
		fmt.Println("parse requestTime failed.")
	}

	log.status, err = record.Field("status")
	if err != nil {
		log.status = ""
		fmt.Println("parse status failed.")
	}

	log.timeLocal, err = record.Field("time_local")
	if err != nil {
		log.timeLocal = ""
		fmt.Println("parse timeLocal failed.")
	}

	log.upstreamAddr, err = record.Field("upstream_addr")
	if err != nil {
		log.upstreamAddr = ""
		fmt.Println("parse upstreamAddr failed.")
	}

	log.upstreamResponseTime, err = record.Field("upstream_response_time")
	if err != nil {
		log.upstreamResponseTime = ""
		fmt.Println("parse upstreamResponseTime failed.")
	}
	return log
}
