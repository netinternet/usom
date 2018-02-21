package Usom

import (
	"context"
	"encoding/xml"
	"net"
	"time"

	"github.com/astaxie/beego/httplib"
)

func Scandaily(masks []string, speed int) map[string]interface{} {
	list := []Pong{}
	scanned := []Pong{}
	var s, _ = time.Parse("2006-01-02 15:04:05", time.Now().Local().Format("2006-01-02 15:04:05"))
	var lastDay = s.Unix() - 86400
	usomUrlList, _ := httplib.Get("https://www.usom.gov.tr/url-list.xml").Bytes()
	xml.Unmarshal(usomUrlList, &t)

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
