package main

import (
	"crypto/tls"
	"fmt"
	"time"
	"flag"
	"strings"
	"encoding/xml"
	"strconv"
)

type result struct {
	Channel string      `xml:"channel"`
	Value string `xml:"value"`
	Unit string `xml:"unit"`
	CustomUnit string `xml:"CustomUnit"`
	Valuelookup string `xml:"ValueLookup"`
}

type prtgbody struct {
	XMLName xml.Name `xml:"prtg"`
	TextField string `xml:"text"`
	Res []result `xml:"result"`
}


type httpcheck struct {
	url string
	error string
	expiredSecs int64
	expiredDays int64
}

func main() {
	alltextmessage := ""
	var urllistassinglestring string
	var urllist []string
	var urlreclist []httpcheck
	urlreclist = make([]httpcheck,0)
	flag.StringVar(&urllistassinglestring,"u","","URL without https prefix, separated by commas")
	flag.Parse()
	if len(urllistassinglestring) > 3{
		urllist = strings.Split(urllistassinglestring,",")
	}

	for _,url := range urllist{
		sterr := ""
		conn, _ := tls.Dial("tcp", url+":443", nil)
		err := conn.VerifyHostname(url)
		expiry := conn.ConnectionState().PeerCertificates[0].NotAfter
		ux := expiry.Unix()
		now := time.Now().Unix()
		uxexpire := ux - now

		if err != nil{
			sterr = fmt.Sprintf("%s",err)
		}
		urlreclist = append(urlreclist,httpcheck{url:url,error:sterr,expiredSecs:uxexpire,expiredDays:uxexpire/86400})
	}
	var rd1 []result
	for _,memberin := range urlreclist{
		ValDate := memberin.expiredDays
		if memberin.error != "" {
			ValDate = -1
			alltextmessage = alltextmessage + ", ERROR:" + memberin.url + " " + memberin.error
		}else {
			if memberin.expiredSecs < 0{
				ValDate =0
			}
		}
		rd1 = append(rd1,result{Channel:memberin.url,Value:strconv.FormatInt(ValDate,10),Unit:"Custom",CustomUnit:"days",Valuelookup:"certcheck"})
	}
	mt1 := &prtgbody{TextField:alltextmessage,Res: rd1}
	bolB, _ := xml.Marshal(mt1)
	fmt.Println(string(bolB))
}
