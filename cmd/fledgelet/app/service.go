// Copyright (c) 2021 Cisco Systems, Inc. and its affiliates
// All rights reserved
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	pbAgent "wwwin-github.cisco.com/eti/fledge/pkg/proto/go/agent"
	"wwwin-github.cisco.com/eti/fledge/pkg/util"
)

const (
	envAgentId = "FLEDGE_AGENT_ID"
)

type AgentService struct {
	apiserverEp string
	notifierEp  string
	name        string
	id          string
	nHandler    *NotifyHandler
}

func NewAgent(apiserverEp string, notifierEp string) (*AgentService, error) {
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	agentId := os.Getenv(envAgentId)
	// in case agent id env variable is not set
	if agentId == "" {
		// TODO: revisit name and id part;
		// determining name and id can be done through api call
		hash := md5.Sum([]byte(name))
		id := hex.EncodeToString(hash[:])

		if name == "" || id == "" {
			err := fmt.Errorf("missing fledgelet name or id")
			zap.S().Error(err)
			return nil, err
		}

		agentId = id
	}

	nHandler := newNotifyHandler(apiserverEp, notifierEp, name, agentId)

	agent := &AgentService{
		apiserverEp: apiserverEp,
		notifierEp:  notifierEp,
		name:        name,
		id:          agentId,
		nHandler:    nHandler,
	}

	return agent, nil
}

func (agent *AgentService) Start() error {
	zap.S().Infof("Starting %s... name: %s | id: %s", util.Agent, agent.name, agent.id)

	agent.nHandler.start()

	err := agent.startAppServer()

	return err
}

// startAppServer starts the fledgelet grpc server and register the corresponding stores implemented by agentServer.
func (agent *AgentService) startAppServer() error {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(util.AgentGrpcPort))
	if err != nil {
		zap.S().Errorf("Failed to listen grpc server: %v", err)
		return err
	}

	// create grpc server
	s := grpc.NewServer()
	server := &appServer{}
	server.init()

	//register grpc services
	pbAgent.RegisterStreamingStoreServer(s, server)

	zap.S().Infof("Agent GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		zap.S().Errorf("Failed to serve: %s", err)
		return err
	}

	return nil
}
