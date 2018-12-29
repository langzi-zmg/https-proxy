package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"regexp"
	"os"
	"bufio"
	"io"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false
	proxy.OnRequest().HandleConnectFunc(
		func(req string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			parts := strings.Split(ctx.Req.RemoteAddr, ":")
			var ip string
			if len(parts) > 0 {
				ip = parts[0]
			} else {
				return goproxy.RejectConnect, req
			}

			if !IPInWhitelist(ip) {
				return goproxy.RejectConnect, req
			}
			return goproxy.OkConnect, req
		})

	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		parts := strings.Split(req.RemoteAddr, ":")
		var ip string
		if len(parts) > 0 {
			ip = parts[0]
		} else {
			return req, nil
		}
		if !IPInWhitelist(ip) {
			return req, goproxy.NewResponse(req,
				goproxy.ContentTypeText, http.StatusForbidden, "")
		}
		return req, nil
	})
	log.Fatal(http.ListenAndServe(":8888", proxy))
}
var IPPool = []string{}
var regs []*regexp.Regexp

func init() {
	fi,err := os.Open("/home/ubuntu/ip")
	if err != nil{
		panic("fail to initialize ip whitelist")
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		ip, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		regs = append(regs, regexp.MustCompile(string(ip)))
	}
}

func IPInWhitelist(given string) bool {
	for _, reg := range regs {
		rs := reg.FindString(given)
		if rs == given {
			return true
		}
	}
	return false
}
