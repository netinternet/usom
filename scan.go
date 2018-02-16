package Usom

import (
	"context"
	"encoding/xml"
	"net"
	"time"

	"github.com/astaxie/beego/httplib"
	"gopkg.in/cheggaaa/pb.v1"
)

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
