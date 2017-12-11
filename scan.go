package usom

import (
	"context"
	"encoding/xml"
	"net"
	"strings"
	"time"

	"github.com/astaxie/beego/httplib"
	"gopkg.in/cheggaaa/pb.v1"
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

/*
clean up the url distorting characters
*/
func Cleanurl(url string) string {
	if strings.Contains(url, "http://") {
		return strings.Split(url, "http://")[1]
	}
	if strings.ContainsAny(url, "/") {
		return strings.Split(url, "/")[0]
	}
	return url
}

/*
the given ip address on the ip mask
return ip address or nil

usage : isInside("192.168.1.5", "{'192.168.1.0/24'}")
*/
func Isinside(adress string, masks []string) net.IP {
	if len(adress) != 0 {
		ipv4 := adress
		p := net.ParseIP(ipv4).To4()
		for i := 0; i < len(masks); i++ {
			_, pc, _ := net.ParseCIDR(masks[i])
			if pc.Contains(p) {
				return p
			}
		}
	}
	return nil
}

/*
Scan usom blocked url list inside your ip mask
if it find , assign list struct slice
masks are contains ip ranges
ex: var masks = []string{"31.192.208.0/21", "89.43.28.0/22"}
the speed sleep down as the speed values increases

usage: Scan(masks, 10)

if view the results

	list := Scan(masks,10)
	for _, v := range list {
		fmt.Println(v.Hostname, v.IP)
	}

*/
func Scan(masks []string, speed int) []Pong {
	list := []Pong{}
	bar := pb.StartNew(0)
	usomUrlList, _ := httplib.Get("https://www.usom.gov.tr/url-list.xml").Bytes()
	xml.Unmarshal(usomUrlList, &t)

	jobs := make(chan string)
	done := make(chan bool)
	go func() {
		for {
			bar.Increment()
			j, more := <-jobs
			if more {
				ctx, _ := context.WithTimeout(context.Background(), time.Duration(speed)*time.Millisecond)
				ipaddr, _ := net.DefaultResolver.LookupHost(ctx, Cleanurl(j))
				if ipaddr != nil {
					for _, v := range ipaddr {
						if Isinside(v, masks) != nil {
							list = append(list, Pong{IP: v, Hostname: Cleanurl(j)})
						}
					}
				}
			} else {
				done <- true
				return
			}
		}
	}()
	for _, url := range t.UrlList.UrlInfo {
		jobs <- url.Url
	}
	close(jobs)
	<-done
	return list
}

func Scandaily(masks []string, speed int) map[string]interface{} {
	list := []Pong{}
	scanned := []Pong{}
	var s, _ = time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	var lastDay = s.Unix() - 86400

	jobs := make(chan string)
	done := make(chan bool)
	go func() {
		for {
			j, more := <-jobs
			if more {
				ctx, _ := context.WithTimeout(context.Background(), time.Duration(speed)*time.Millisecond)
				ipaddr, _ := net.DefaultResolver.LookupHost(ctx, Cleanurl(j))
				if ipaddr != nil {
					for _, v := range ipaddr {
						scanned = append(scanned, Pong{IP: v, Hostname: Cleanurl(j)})
						if Isinside(v, masks) != nil {
							list = append(list, Pong{IP: v, Hostname: Cleanurl(j)})
						}
					}
				}
			} else {
				done <- true
				return
			}
		}
	}()
	for _, url := range t.UrlList.UrlInfo {
		p, _ := time.Parse("2006-01-02 15:04:05", url.Date)
		if p.Unix() > lastDay {
			jobs <- url.Url
		}
	}
	close(jobs)
	<-done

	e := make(map[string]interface{})
	e["scanned"] = scanned
	e["results"] = list
	return e
}
