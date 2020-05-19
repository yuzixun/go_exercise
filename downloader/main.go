package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	yaml "gopkg.in/yaml.v2"
)

type Conf struct {
	Database struct {
		Name string `yaml:"name"`
		User string `yaml:"user"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Pwd  string `yaml:"pwd"`
	} `yaml:"database"`
}

type Fetcher struct {
	url      string
	fileName string
}

func main() {
	file, err := ioutil.ReadFile("conf/db.yml")
	if err != nil {
		log.Fatal("can not open file", err)
	}

	var config Conf
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("can't parse config file ", err)
	}

	configStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", config.Database.User, config.Database.Pwd, config.Database.Host, config.Database.Port, config.Database.Name)
	db, err := gorm.Open("mysql", configStr)
	if err != nil {
		log.Fatal("gorm open ", err)
	}
	defer db.Close()

	modelChan := make(chan CallRecord, 10000)
	quitChan := make(chan int)
	fetchChan := make(chan Fetcher, 10000)

	go controller(modelChan, quitChan, fetchChan)
	go queryCallRecords(modelChan, quitChan, db)

	for i := 0; i < 20; i++ {
		go handleDownload(fetchChan)
	}

	time.Sleep(1000 * time.Second)
}

func controller(modelChan chan CallRecord, quitChan chan int, fetchChan chan Fetcher) {
	for {
		select {
		case callRecord := <-modelChan:
			fileProcessor(callRecord, fetchChan)
		case <-quitChan:
			fmt.Println("quit")
			return
		}
	}
}

func fileProcessor(callRecord CallRecord, fetchChan chan Fetcher) {
	if len(callRecord.FileName) == 0 || callRecord.FileName == "NULL" || !strings.HasPrefix(callRecord.FileName, "https://") {
		return
	}

	path := fmt.Sprintf("export/%d/%d/%d/%d", callRecord.OrganizationId, callRecord.CreatedAt.Year(), callRecord.CreatedAt.Month(), callRecord.CreatedAt.Day())
	os.MkdirAll(path, os.ModePerm)

	strs := strings.Split(callRecord.FileName, "/")
	fileName := fmt.Sprintf("%s/%s", path, strs[len(strs)-1])

	fetchChan <- Fetcher{url: callRecord.FileName, fileName: fileName}
}
