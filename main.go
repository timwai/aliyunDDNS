package main

import (
	"flag"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net/http"
)

var accessKeyId = flag.String("id", "", "")
var accessKeySecret = flag.String("secret", "", "")
var domain = flag.String("domain", "", "")
var RR = flag.String("RR", "", "")
var recordType = flag.String("type", "", "")

func main() {
	flag.Parse()
	client, _ := alidns.NewClientWithAccessKey("cn-hangzhou", *accessKeyId, *accessKeySecret)
	record := getDomainRecords(client)
	ip := getIP()
	if record == nil {
		addDomain(client, ip)
	} else {
		value := record.Value
		recordId := record.RecordId
		if value == ip {
			log.Println("相关IP地址记录已存在，无需更新解析记录")
			return
		}
		updateDomain(client, ip, recordId)
	}
}

func addDomain(client *alidns.Client, ip string) {
	request := alidns.CreateAddDomainRecordRequest()
	request.Scheme = "https"
	request.Value = ip
	request.Type = *recordType
	request.RR = *RR
	request.DomainName = *domain
	response, err := client.AddDomainRecord(request)
	if err != nil {
		fmt.Print(err.Error())
		log.Println("新增解析记录失败")
	}
	if response.RecordId != "" {
		log.Println("新增解析记录成功")
	} else {
		log.Println("新增解析记录失败")
	}
}

func updateDomain(client *alidns.Client, ip string, recordId string) {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.Value = ip
	request.Type = *recordType
	request.RR = *RR
	request.RecordId = recordId
	response, err := client.UpdateDomainRecord(request)
	if err != nil {
		fmt.Print(err.Error())
		log.Println("更新解析记录失败")
	}
	if response.RecordId != "" {
		log.Println("更新解析记录成功")
	} else {
		log.Println("更新解析记录失败")
	}
}

func getDomainRecords(client *alidns.Client) *alidns.Record {
	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"
	request.Type = *recordType
	request.DomainName = *domain
	request.SearchMode = "ADVANCED"
	request.RRKeyWord = *RR
	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	count := response.TotalCount
	if count == 0 {
		return nil
	}
	records := response.DomainRecords.Record
	return &records[0]
}

func getIP() string {
	if *recordType == "AAAA" {
		return getIPV6()
	} else {
		return getIPV4()
	}
}

func getIPV4() string {
	//https://ipv4.netarm.com/
	resp, err := http.Get("https://4.ipw.cn")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}

func getIPV6() string {
	//https://ipv6.netarm.com/
	resp, err := http.Get("https://6.ipw.cn")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	return string(content)
}
