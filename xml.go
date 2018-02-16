package Usom

import (
	"encoding/xml"

	"github.com/astaxie/beego/httplib"
)

/*
XML Parsing Structs for https://www.usom.gov.tr/url-list.xml
*/
type UsomData struct {
	XMLName xml.Name `xml:"usom-data"`
	XMLInfo XMLInfo  `xml:"xml-info"`
	UrlList UrlList  `xml:"url-list"`
}

type XMLInfo struct {
	XMLName xml.Name `xml:"xml-info"`
	Updated string   `xml:"updated"`
	Author  string   `xml:"author"`
}

type UrlList struct {
	XMLName xml.Name  `xml:"url-list"`
	UrlInfo []UrlInfo `xml:"url-info"`
}

type UrlInfo struct {
	XMLName xml.Name `xml:"url-info"`
	Id      int32    `xml:"id"`
	Url     string   `xml:"url"`
	Desc    string   `xml:"desc"`
	Source  string   `xml:"source"`
	Date    string   `xml:"date"`
}

/*
Pong structs
for the ip address and hostnames in the given network
shortly: result struct
*/

type Pong struct {
	IP       string
	Hostname string
}

/*
Global Variables
*/
var t = UsomData{}
var usomUrlList, _ = httplib.Get("https://www.usom.gov.tr/url-list.xml").Bytes()
var done = xml.Unmarshal(usomUrlList, &t)
