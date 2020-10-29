/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package p4rt implements p4runtime functions
*/
package p4rt

import v1 "github.com/p4lang/p4runtime/go/p4/v1"

type readChan struct {
	client       v1.P4Runtime_ReadClient
	responseChan chan *v1.ReadResponse
}

func (r readChan) Recv() {
	for {
		log.Debug("In Recv for loop")
		readResp, err := r.client.Recv()
		log.Debug("In Recv for loop after receiving message")
		if err != nil {
			log.Debugf("Failed to receive a message: %v\n", err)
			return
		}
		r.responseChan <- readResp
	}
}
