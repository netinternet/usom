package scanner

import (
	"encoding/xml"
	"errors"
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

/*
Lookup Method with timeout duration
works concurrently

usage: Lookup("garantinternetsubem.com", time.Millisecond*100)
note: If "http://" or "/asd/qwe" are in the url , function it will not work properly
therefore you should use with cleanUrl method
usage: Lookup(cleanUrl("garantinternetsubem.com"), time.Millisecond*100)
*/
func Lookup(hostname string, timeout time.Duration) ([]string, error) {
	c1 := make(chan []string)
	c2 := make(chan error)

	var ipaddr []string
	var err error

	go func() {
		var ipaddr []string
		ipaddr, err := net.LookupHost(hostname)
		if err != nil {
			c2 <- err
		}
		c1 <- ipaddr
	}()

	select {
	case ipaddr = <-c1:
	case err = <-c2:
	case <-time.After(timeout):
		return ipaddr, errors.New("Timeout")
	}
	if err != nil {
		return ipaddr, errors.New("Timeout")
	}
	return ipaddr, nil
}

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

value is total url search count

filtertime is value of time filtering as timestamp


usage: usom(masks, time.Millisecond*100)

if view the results

	list := usom(masks,time.Millisecond*100)
	for _, v := range list {
		fmt.Println(v.Hostname, v.IP)
	}

*/
func Usom(masks []string, speed time.Duration) []Pong {
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
				ipaddr, _ := Lookup(Cleanurl(j), speed)
				if ipaddr != nil {
					for _, v := range ipaddr {
						if Isinside(v, masks) != nil {
							temp := Pong{IP: v, Hostname: Cleanurl(j)}
							list = append(list, temp)
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

func UsomDaily(masks []string, speed time.Duration) []Pong {
	list := []Pong{}
	var s, _ = time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	var lastDay = s.Unix() - 86400
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
				ipaddr, _ := Lookup(Cleanurl(j), speed)
				if ipaddr != nil {
					for _, v := range ipaddr {
						if Isinside(v, masks) != nil {
							temp := Pong{IP: v, Hostname: Cleanurl(j)}
							list = append(list, temp)
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
	return list
}
