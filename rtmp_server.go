//
// Copyright 2014-2099 Hong Miao. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import(
	"fmt"
	"flag"
	"github.com/golang/glog"
	"github.com/oikomi/gortmpserver/server"
	"github.com/oikomi/gortmpserver/libnet"
)

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
const char* build_time(void) {
	static const char* psz_build_time = "["__DATE__ " " __TIME__ "]";
	return psz_build_time;
}
*/
import "C"

var (
	buildTime = C.GoString(C.build_time())
)

func BuildTime() string {
	return buildTime
}

const VERSION string = "0.10"

func version() {
	fmt.Printf("rtmp_server version %s Copyright (c) 2014-2099 Harold Miao (miaohonghit@gmail.com)  \n", VERSION)
}

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

var InputConfFile = flag.String("conf_file", "./server/server.json", "input conf file name") 

func main() {
	version()
	fmt.Printf("built on %s\n", BuildTime())
	flag.Parse()
	cfg := server.NewRtmpServerConfig(*InputConfFile)
	err := cfg.LoadConfig()
	if err != nil {
		glog.Error(err.Error())
		return
	}
	
	s, err := libnet.Listen(cfg.TransportProtocols, cfg.Listen)
	if err != nil {
		glog.Error(err.Error())
		return
	}
	
	glog.Info("rtmp server start at ", s.Listener().Addr().String())
	rs := server.NewRtmpServer(cfg, s)
	
	rs.ServerLoop()
}
