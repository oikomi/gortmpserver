//
// Copyright 2014 Hong Miao. All Rights Reserved.
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

import (
	"github.com/oikomi/gortmpserver/server/config"
	"github.com/oikomi/gortmpserver/server/rtmpserver"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(0)
	}
	
	cfg, err := config.LoadConfig(os.Args[1])
	log.Println(cfg.Listen)
	if err != nil {
		log.Fatalln(err.Error())
		return 
	}
	
	rtmpserver := rtmpserver.NewRtmpServer(cfg)
	err = rtmpserver.Listen()
	if err != nil {
		log.Fatalln(err.Error())
		return 
	}
	
	
	
}
