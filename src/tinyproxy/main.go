package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"
	//"fmt"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"

	"github.com/go-ini/ini"
)

var (
	verbose = flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr    = flag.String("addr", ":3128", "proxy listen address")
)

func init() {
	flag.Parse()
}

func main() {
	var cfgAuth *ini.Section
	if cfg, err := ini.Load("conf/tinyproxy.conf"); err == nil {
		cfgAuth, _ = cfg.GetSection("auth")
	}
	//fmt.Println(cfgAuth)

	proxy := goproxy.NewProxyHttpServer()
	auth.ProxyBasic(proxy, "hello", func(user, passwd string) bool {
		if userKey, err := cfgAuth.GetKey(user); err != nil {
			return false
		} else {
			//fmt.Println(user+" "+userKey.String())
			return userKey.String() == passwd
		}
	})
	proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*baidu.com$"))).
		HandleConnect(goproxy.AlwaysReject)
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
