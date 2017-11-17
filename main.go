/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"log"
	"os"
	"runtime"
	"strings"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

var nc *nats.Conn
var cfg *ecc.Config

func handler(msg *nats.Msg) {
	var e Event
	var err error

	log.Println(msg.Subject)

	defer response(msg.Subject, &e, &err)

	e, err = event(msg.Subject, msg.Data)
	if err != nil {
		return
	}

	parts := strings.Split(msg.Subject, ".")

	switch parts[1] {
	case "create":
		err = e.Create()
	case "update":
		err = e.Update()
	case "delete":
		err = e.Delete()
	case "find":
		err = e.Find()
	default:
		err = errors.New("unsupported action type")
	}
}

func main() {
	cfg = ecc.NewConfig(os.Getenv("NATS_URI"))
	nc = cfg.Nats()

	nc.Subscribe("*.*.vcloud", handler)

	runtime.Goexit()
}
