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

package rtmpserver

import (
	"errors"
	"net/url"
	"strings"
)

type Request struct {
	TcUrl string
	PageUrl string
	SwfUrl string
	ObjectEncoding float64

	Schema string
	Vhost string
	Port string
	App string
	Stream string
}

func NewRequest() (*Request) {
	return &Request {
		
	}
}


func (self *Request) DiscoveryApp() error {
	var err error
	var v string = self.TcUrl
	if !strings.Contains(v, "?") {
		v = strings.Replace(v, "...", "?", 1)
		v = strings.Replace(v, "...", "=", 1)
	}
	for strings.Contains(v, "...") {
		v = strings.Replace(v, "...", "&", 1)
		v = strings.Replace(v, "...", "=", 1)
	}
	self.TcUrl = v

	var u *url.URL
	if u, err = url.Parse(self.TcUrl); err != nil {
		return err
	}

	self.Schema, self.App = u.Scheme, u.Path

	self.Vhost = u.Host
	if strings.Contains(u.Host, ":") {
		host_parts := strings.Split(u.Host, ":")
		self.Vhost, self.Port = host_parts[0], host_parts[1]
	}

	// discovery vhost from query.
	query := u.Query()
	for k, _ := range query {
		if strings.ToLower(k) == "vhost" && query.Get(k) != "" {
			self.Vhost = query.Get(k)
		}
	}

	// resolve the vhost from config
	// TODO: FIXME: implements it
	// TODO: discovery the params of vhost.

	if self.Schema = strings.Trim(self.Schema, "/\n\r "); self.Schema == ""{
		return errors.New("discovery schema failed")
	}
	if self.Vhost = strings.Trim(self.Vhost, "/\n\r "); self.Vhost == "" {
		return errors.New("discovery vhost failed")
	}
	if self.App = strings.Trim(self.App, "/\n\r "); self.App == "" {
		return errors.New("discovery app failed. tcUrl")
	}
	if self.Port = strings.Trim(self.Port, "/\n\r "); self.Port == "" {
		return errors.New("discovery port failed")
	}

	return nil
}