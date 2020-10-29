/*
 * Copyright 2019-present Open Networking Foundation
 *
 * SPDX-License-Identifier: Apache-2.0
 */

/*
Package p4rt implements p4runtime functions
*/
package p4rt

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	v1 "github.com/p4lang/p4runtime/go/p4/v1"
	"github.com/stratum/testvectors-runner/pkg/utils/common"
	tvb "github.com/stratum/testvectors/proto/target"
	"google.golang.org/grpc"
)

//Connection struct stores the P4Runtime client connection, context and cancel function.
type connection struct {
	ctx       context.Context
	client    v1.P4RuntimeClient
	connError error
	cancel    context.CancelFunc
}

//connect starts a gRPC connection to the target specified.
//It returns connection struct with P4Runtime client, close function
//If an error is encountered during opening the connection, it is returned.
func connect(tg *tvb.Target) connection {
	log.Debug("In p4_oper connect")
	if tg.Address == "" {
		return connection{connError: errors.New("an address must be specified")}
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, CtxTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, tg.Address, grpc.WithInsecure())
	if err != nil {
		return connection{connError: fmt.Errorf("cannot dial target %s, %v", tg.Address, err)}
	}
	return connection{ctx: ctx, client: v1.NewP4RuntimeClient(conn), cancel: func() { conn.Close() }}
}

//Write calls P4RuntimeClient's Write and returns WriteResponse
func (c connection) Write(writeReq *v1.WriteRequest) *v1.WriteResponse {
	log.Info("Sending P4 write request")
	log.Debugf("Write request: %s", writeReq)
	ctx := context.Background()
	resp, err := c.client.Write(ctx, writeReq)
	if err != nil {
		log.Errorf("Error sending P4 write request:%v", err)
		return nil
	}
	log.Infof("Received P4 write response")
	log.Debugf("P4 write response:%s", resp)
	return resp
}

//Read calls P4RuntimeClient's Read and returns ??
func (c connection) Read(readReq *v1.ReadRequest) readChan {
	log.Info("Sending P4 read request")
	log.Debugf("Read request: %s", readReq)
	ctx := context.Background()
	readClient, err := c.client.Read(ctx, readReq)
	if err != nil {
		log.Errorf("Error sending P4 read request:%v", err)
		return readChan{}
	}
	return readChan{client: readClient, responseChan: make(chan *v1.ReadResponse)}
}

//SetForwardingPipelineConfig calls P4RuntimeClient's SetForwardingPipelineConfig and returns SetForwardingPipelineConfigResponse
func (c connection) SetForwardingPipelineConfig(pipelineCfg *v1.SetForwardingPipelineConfigRequest) *v1.SetForwardingPipelineConfigResponse {
	log.Info("Sending P4 pipeline config")
	log.Debugf("Pipeline config: %s", pipelineCfg)
	ctx := context.Background()
	resp, err := c.client.SetForwardingPipelineConfig(ctx, pipelineCfg)
	if err != nil {
		log.Errorf("Error sending P4 pipeline config:%v", err)
		return nil
	}
	log.Info("Received P4 pipeline config response")
	log.Debugf("P4 set pipeline config response:%s\n", resp)
	return resp
}

//GetForwardingPipelineConfig calls P4RuntimeClient's GetForwardingPipelineConfig and returns GetForwardingPipelineConfigResponse
func (c connection) GetForwardingPipelineConfig(pipelineCfg *v1.GetForwardingPipelineConfigRequest) *v1.GetForwardingPipelineConfigResponse {
	log.Info("Getting P4 pipeline config")
	log.Debugf("Get Pipeline config: %s", pipelineCfg)
	ctx := context.Background()
	resp, err := c.client.GetForwardingPipelineConfig(ctx, pipelineCfg)
	if err != nil {
		log.Errorf("Error getting P4 pipeline config:%v", err)
		return nil
	}
	log.Info("Received get P4 pipeline config response")
	log.Debugf("P4 get pipeline config response:%s\n", resp)
	return resp
}

//verifyWriteResp compares two WriteResponses and returns true or false
func verifyWriteResp(expected, actual *v1.WriteResponse) bool {
	//FIXME
	//initializing expected to empty response to avoid nil pointer exception when tv doesn't have response
	if expected == nil {
		expected = &v1.WriteResponse{}
	}
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both write responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Write responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case proto.Equal(expected, actual):
		log.Info("Write responses are equal")
		log.Debugf("Write response: %s\n", actual)
		return true
	default:
		log.Warnf("Write responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	}
}

func verifyReadRespList(expResp []*v1.ReadResponse, actRespChan chan *v1.ReadResponse) bool {
	result := true
	for _, exp := range expResp {
		select {
		case act := <-actRespChan:
			log.Debug("In verifyReadRespList Case read response channel")
			result = result && verifyReadResp(exp, act)
		case <-time.After(PktTimeout):
			log.Error("Timed out waiting for read response")
			return false
		}
	}
	return result
}

//verifyReadResp compares two ReadResponses and returns true or false
func verifyReadResp(expected, actual *v1.ReadResponse) bool {
	//FIXME
	//initializing "expected" to empty read response to avoid nil pointer exception when tv doesn't have response
	if expected == nil {
		expected = &v1.ReadResponse{}
	}
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both read responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Read responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case proto.Equal(expected, actual):
		log.Info("Read responses are equal")
		log.Debugf("Read response: %s\n", actual)
		return true
	default:
		log.Warnf("Read responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	}
}

//verifySetForwardingPipelineConfigResp compares two SetForwardingPipelineConfigResponse and returns true or false
func verifySetForwardingPipelineConfigResp(expected, actual *v1.SetForwardingPipelineConfigResponse) bool {
	//FIXME
	//initializing expected to empty response to avoid nil pointer exception when tv doesn't have response
	if expected == nil {
		expected = &v1.SetForwardingPipelineConfigResponse{}
	}

	switch {
	case expected == nil && actual == nil:
		log.Debug("Both SetForwardingPipelineConfig responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("SetForwardingPipelineConfig responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case proto.Equal(expected, actual):
		log.Info("SetForwardingPipelineConfig responses are equal")
		log.Debugf("SetForwardingPipelineConfig response: %s\n", actual)
		return true
	default:
		log.Warnf("SetForwardingPipelineConfig responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	}
}

//verifyGetForwardingPipelineConfigResp compares two GetForwardingPipelineConfigResponse and returns true or false
func verifyGetForwardingPipelineConfigResp(expected, actual *v1.GetForwardingPipelineConfigResponse) bool {
	//FIXME
	//initializing expected to empty response to avoid nil pointer exception when tv doesn't have response
	if expected == nil {
		expected = &v1.GetForwardingPipelineConfigResponse{}
	}

	switch {
	case expected == nil && actual == nil:
		log.Debug("Both GetForwardingPipelineConfig responses are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("GetForwardingPipelineConfig responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	case proto.Equal(expected, actual):
		log.Info("GetForwardingPipelineConfig responses are equal")
		log.Debugf("GetForwardingPipelineConfig response: %s\n", actual)
		return true
	default:
		log.Warnf("GetForwardingPipelineConfig responses are unequal\nExpected: %s\nActual  : %s\n", expected, actual)
		return false
	}
}

//verifyPacketIn compares two PacketIns and returns true or false
func verifyPacketIn(expected, actual *v1.PacketIn) bool {
	switch {
	case expected == nil && actual == nil:
		log.Debug("Both packets are empty")
		return true
	case expected == nil || actual == nil:
		log.Warnf("Packets don't match\nExpected: % s\nActual  : % s\n", expected, actual)
		return false
	case !bytes.Equal(expected.GetPayload(), actual.GetPayload()):
		log.Warnf("Payloads don't match\nExpected: % x\nActual  : % x\n", expected.GetPayload(), actual.GetPayload())
		return false
	case !compareMetadata(expected.GetMetadata(), actual.GetMetadata()):
		log.Warnf("Metadata don't match\nExpected: % v\nActual  : % v\n", expected.GetMetadata(), actual.GetMetadata())
		return false
	default:
		log.Info("PacketIns are equal")
		log.Debugf("PacketIn: %s", actual)
		return true
	}
}

func compareMetadata(m1, m2 []*v1.PacketMetadata) bool {
	switch {
	case len(m1) < 1 && len(m2) < 1:
		log.Debug("Both metadata are empty")
		return true
	case len(m1) < 1 || len(m2) < 1:
		log.Debug("Metadata don't match")
		return false
	default:
		//TODO: verify complete metadata instead of only one
		//return reflect.DeepEqual(m1, m2)
		/*for i := 0; i < len(m1); i++ {
			if m1[i].GetMetadataId() != m2[i].GetMetadataId() || common.GetInt(m1[i].GetValue()) != common.GetInt(m2[i].GetValue()) {
				return false
			}
		}
		return true*/
		return common.GetInt(m1[0].Value) == common.GetInt(m2[0].Value)
	}
}
