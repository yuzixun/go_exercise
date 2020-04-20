package main

import (
	"fmt"
	"io"
	"os"

	"github.com/satyrius/gonx"
)

func parseLine(line string) *LogData {
	var logReader io.Reader
	var err error
	file, err := os.Open("./access.log")
	if err != nil {
		panic(err)
	}
	logReader = file
	defer file.Close()

	// logReader = strings.NewReader(line)
	// fmt.Println(logReader)
	format := "$remote_addr  $remote_user [$time_local] \"$request\" \"$http_host\" $status $request_length $body_bytes_sent \"$http_referer\" \"$http_user_agent\"     $http_x_forwarded_for   \"$upstream_addr\" \"$request_time\"        \"$upstream_response_time\""

	reader := gonx.NewReader(logReader, format)
	// fmt.Printf(" %+v\n", reader)
	for {
		rec, err := reader.Read()
		fmt.Printf("Parsed entry: %+v\n", rec)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		// Process the record... e.g.
	}

	// line = strings.TrimSpace(line)
	// fmt.Println("line si ", line)
	// // nginxReg := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	// nginxReg := regexp.MustCompile(`?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}) - - \[(?P<dateandtime>\d{2}\/[a-z]{3}\/\d{4}:\d{2}:\d{2}:\d{2} (\+|\-)\d{4})\] ((\"(GET|POST) )(?P<url>.+)(http\/1\.1")) (?P<statuscode>\d{3}) (?P<bytessent>\d+) (["](?P<refferer>(\-)|(.+))["]) (["](?P<useragent>.+)["]`)
	// // nginxReg := regexp.MustCompile(`(?\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})?(?\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})?-? - ?\S* \[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (\+|\-)\d{4})\]\s+\"(?P<method>\S{3,10}) (?P<path>\S+) HTTP\/1\.\d" (?P<response_status>\d{3}) (?P<bytes>\d+) "(?P<referer>(\-)|(.+))?" "(?P<useragent>.+)`)
	// strs := nginxReg.FindAllStringSubmatch(line, -1)

	// fmt.Println("strs is", strs)

	// // data := strings.FieldsFunc(line, Split)
	// // substrings := strings.SplitN(line, " ", 4)
	// // fmt.Printf("%v\n", substrings)
	// // fmt.Printf("%#v\n", substrings)
	// // fmt.Printf("%+v\n", substrings)
	// // remoteAddr, remoteUser, line := substrings[0], substrings[2], substrings[3]
	// // line = substrings[3]
	// // line.SplitN()
	// // regexp.MustComplie("[*]+")
	return &LogData{
		remoteAddr:           "remoteAddr",
		remoteUser:           "remoteUser",
		timeLocal:            "data[2]",
		request:              "data[3]",
		host:                 "data[4]",
		status:               "data[5]",
		requestLength:        "data[6]",
		bodyBytesSent:        "data[7]",
		httpReferer:          "data[8]",
		httpUserAgent:        "data[9]",
		forwardedFor:         "data[10]",
		upstreamAddr:         "data[11]",
		requestTime:          "data[12]",
		upstreamResponseTime: "data[13]",
	}
}
