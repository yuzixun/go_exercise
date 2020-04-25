package parser

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/satyrius/gonx"
)

type NgxLog struct {
	RemoteAddr           string
	RemoteUser           string
	TimeLocal            string
	Time                 time.Time
	Request              string
	RequestURL           string
	HttpHost             string
	Status               string
	RequestLength        string
	BodyBytesSent        string
	HttpReferer          string
	HttpUserAgent        string
	ForwardedFor         string
	UpstreamAddr         string
	RequestTime          string
	UpstreamResponseTime string
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
		// fmt.Printf("%+v\n", rec)

		logProcessor <- generateNgxLogBy(rec)
		// time.Sleep(time.Millisecond * 1)
		// break
	}
}

func generateNgxLogBy(record *gonx.Entry) NgxLog {
	var err error
	log := NgxLog{}

	log.BodyBytesSent, err = record.Field("body_bytes_sent")
	if err != nil {
		log.BodyBytesSent = ""
		fmt.Println("parse BodyBytesSent failed.")
	}

	log.HttpHost, err = record.Field("http_host")
	if err != nil {
		log.HttpHost = ""
		fmt.Println("parse HttpHost failed.")
	}

	log.HttpReferer, err = record.Field("http_referer")
	if err != nil {
		log.HttpReferer = ""
		fmt.Println("parse HttpReferer failed.")
	}

	log.HttpUserAgent, err = record.Field("http_user_agent")
	if err != nil {
		log.HttpUserAgent = ""
		fmt.Println("parse HttpUserAgent failed.")
	}

	log.ForwardedFor, err = record.Field("http_x_forwarded_for")
	if err != nil {
		log.ForwardedFor = ""
		fmt.Println("parse ForwardedFor failed.")
	}

	log.RemoteAddr, err = record.Field("remote_addr")
	if err != nil {
		log.RemoteAddr = ""
		fmt.Println("parse RemoteAddr failed.")
	}

	log.RemoteUser, err = record.Field("remote_user")
	if err != nil {
		log.RemoteUser = ""
		fmt.Println("parse RemoteUser failed.")
	}

	log.Request, err = record.Field("request")
	if err != nil {
		log.Request = ""
		fmt.Println("parse Request failed.")
	}

	urls := strings.SplitN(log.Request, "?", 2)
	log.RequestURL = urls[0]

	log.RequestLength, err = record.Field("request_length")
	if err != nil {
		log.RequestLength = ""
		fmt.Println("parse RequestLength failed.")
	}

	log.RequestTime, err = record.Field("request_time")
	if err != nil {
		log.RequestTime = ""
		fmt.Println("parse RequestTime failed.")
	}

	log.Status, err = record.Field("status")
	if err != nil {
		log.Status = ""
		fmt.Println("parse Status failed.")
	}

	log.TimeLocal, err = record.Field("time_local")
	if err != nil {
		log.TimeLocal = ""
		fmt.Println("parse TimeLocal failed.")
	}

	// 13/Apr/2020:03:33:31 +0800
	// 02/Jan/2006:15:04:05 +0700
	log.Time, err = time.Parse("02/Jan/2006:15:04:05 -0700", log.TimeLocal)
	if err != nil {
		log.Time = time.Time{}
		fmt.Println("time err", err)
	}

	log.UpstreamAddr, err = record.Field("upstream_addr")
	if err != nil {
		log.UpstreamAddr = ""
		fmt.Println("parse UpstreamAddr failed.")
	}

	log.UpstreamResponseTime, err = record.Field("upstream_response_time")
	if err != nil {
		log.UpstreamResponseTime = ""
		fmt.Println("parse UpstreamResponseTime failed.")
	}

	return log
}
