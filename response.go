/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

func response(subj string, e *Event, err *error) {
	if *err != nil {
		log.Println(subj + ": " + (*err).Error())
		subj = subj + ".error"

		(*e).SetError(*err)
		(*e).SetState("errored")
	} else {
		subj = subj + ".done"

		(*e).SetState("completed")
	}

	log.Println(subj)

	data, merr := json.Marshal(*e)
	if merr != nil {
		log.Println(merr)
	}

	if err := nc.Publish(subj, data); err != nil {
		log.Println(err)
	}
}
