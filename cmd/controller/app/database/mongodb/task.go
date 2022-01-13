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

package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"github.com/cisco/fledge/cmd/controller/app/objects"
	"github.com/cisco/fledge/pkg/openapi"
	"github.com/cisco/fledge/pkg/util"
)

// CreateTasks creates task records in task db collection
func (db *MongoService) CreateTasks(tasks []objects.Task, dirty bool) error {
	zap.S().Debugf("Calling CreateTasks")

	success := false

	// rollback closure in case of error
	defer func() {
		if success {
			return
		}
		// TODO: implement this
	}()

	for _, task := range tasks {
		cfgData, err := json.Marshal(&task.JobConfig)
		if err != nil {
			return fmt.Errorf("failed to marshal task: %v", err)
		}

		zap.S().Debugf("Creating task for agent %s", task.AgentId)
		filter := bson.M{util.DBFieldJobId: task.JobId, util.DBFieldAgentId: task.AgentId}
		update := bson.M{
			"$set": bson.M{
				util.DBFieldRole:      task.Role,
				util.DBFieldTaskType:  task.Type,
				"config":              cfgData,
				"code":                task.ZippedCode,
				util.DBFieldTaskDirty: dirty,
				util.DBFieldState:     openapi.READY,
				util.DBFieldTimestamp: time.Now(),
			},
		}

		after := options.After
		upsert := true
		opts := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
			Upsert:         &upsert,
		}

		var updatedDoc bson.M
		err = db.taskCollection.FindOneAndUpdate(context.TODO(), filter, update, &opts).Decode(&updatedDoc)
		if err != nil {
			zap.S().Errorf("Failed to insert a task: %v", err)
			return err
		}
	}

	success = true
	zap.S().Debugf("Created tasks successufully")

	return nil
}

func (db *MongoService) GetTask(jobId string, agentId string) (map[string][]byte, error) {
	zap.S().Infof("Fetching task - jobId: %s, agentId: %s", jobId, agentId)
	filter := bson.M{util.DBFieldJobId: jobId, util.DBFieldAgentId: agentId}

	type taskInMongo struct {
		Config []byte `json:"config"`
		Code   []byte `json:"code"`
	}

	var task taskInMongo
	err := db.taskCollection.FindOne(context.TODO(), filter).Decode(&task)
	if err != nil {
		zap.S().Warnf("Failed to fetch task: %v", err)
		return nil, err
	}

	taskMap := map[string][]byte{
		util.TaskConfigFile: task.Config,
		util.TaskCodeFile:   task.Code,
	}

	return taskMap, nil
}

func (db *MongoService) DeleteTasks(jobId string, dirty bool) error {
	zap.S().Infof("Deleting tasks for job: %s", jobId)
	filter := bson.M{util.DBFieldJobId: jobId, util.DBFieldTaskDirty: dirty}

	_, err := db.taskCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		zap.S().Warnf("Failed to remove task: %v", err)
		return err
	}

	return nil
}

func (db *MongoService) UpdateTaskStatus(jobId string, agentId string, taskStatus openapi.TaskStatus) error {
	switch taskStatus.State {
	case openapi.RUNNING:
		fallthrough
	case openapi.FAILED:
		fallthrough
	case openapi.TERMINATED:
		fallthrough
	case openapi.COMPLETED:

	case openapi.READY:
		fallthrough
	case openapi.STARTING:
		fallthrough
	case openapi.APPLYING:
		fallthrough
	case openapi.DEPLOYING:
		fallthrough
	case openapi.STOPPING:
		return fmt.Errorf("prohibited state: %s", taskStatus.State)

	default:
		return fmt.Errorf("unknown state: %s", taskStatus.State)
	}

	filter := bson.M{util.DBFieldJobId: jobId, util.DBFieldAgentId: agentId}

	setElements := bson.M{util.DBFieldState: taskStatus.State, util.DBFieldTimestamp: time.Now()}

	update := bson.M{"$set": setElements}

	after := options.After
	upsert := true
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	updatedDoc := openapi.TaskStatus{}
	err := db.taskCollection.FindOneAndUpdate(context.TODO(), filter, update, &opts).Decode(&updatedDoc)
	if err != nil {
		return ErrorCheck(err)
	}

	return nil
}

// IsOneTaskInState returns true if the state of at least one task is set to the given state.
// Otherwise, it returns false.
func (db *MongoService) IsOneTaskInState(jobId string, state openapi.JobState) bool {
	filter := bson.M{util.DBFieldJobId: jobId, util.DBFieldState: state}

	result := db.taskCollection.FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return false
	}

	return true
}

func (db *MongoService) IsOneTaskInStateWithRole(jobId string, state openapi.JobState, role string) bool {
	filter := bson.M{util.DBFieldJobId: jobId, util.DBFieldState: state, util.DBFieldRole: role}

	result := db.taskCollection.FindOne(context.TODO(), filter)

	if result.Err() != nil {
		return false
	}

	return true
}

func (db *MongoService) SetTaskDirtyFlag(jobId string, dirty bool) error {
	zap.S().Infof("Setting dirty flag to %s tasks for job: %s", dirty, jobId)

	filter := bson.M{util.DBFieldJobId: jobId}
	update := bson.M{
		"$set": bson.M{
			util.DBFieldTaskDirty: dirty,
		},
	}

	_, err := db.taskCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		zap.S().Warnf("Failed to set dirty flag to %s: %v", dirty, err)
		return err
	}

	return nil
}
