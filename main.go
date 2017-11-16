/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
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

func response(subj string, e Event, err *error) {
	if *err != nil {
		log.Println(subj + ": " + (*err).Error())
		subj = subj + ".error"

		e.SetError(*err)
		e.SetState("error")
	} else {
		subj = subj + ".done"

		e.SetState("completed")
	}

	log.Println(subj)

	data, merr := json.Marshal(e)
	if merr != nil {
		log.Println(merr)
	}

	nc.Publish(subj, data)
}

func handler(msg *nats.Msg) {
	var e Event
	var err error

	log.Println(msg.Subject)

	defer response(msg.Subject, e, &err)

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
