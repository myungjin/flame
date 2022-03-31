// Copyright 2022 Cisco Systems, Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

/*
 * Flame REST API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/cisco-open/flame/cmd/controller/app/database"
	"github.com/cisco-open/flame/cmd/controller/app/job"
	"github.com/cisco-open/flame/pkg/openapi"
)

const (
	defaultWaitTime = 10 * time.Second
)

// JobsApiService is a service that implents the logic for the JobsApiServicer
// This service should implement the business logic for every endpoint for the JobsApi API.
// Include any external packages or services that will be required by this service.
type JobsApiService struct {
	dbService database.DBService

	jobEventQ  *job.EventQ
	jobBuilder *job.JobBuilder
}

// NewJobsApiService creates a default api service
func NewJobsApiService(dbService database.DBService, jobEventQ *job.EventQ, jobBuilder *job.JobBuilder) openapi.JobsApiServicer {
	return &JobsApiService{
		dbService: dbService,

		jobEventQ:  jobEventQ,
		jobBuilder: jobBuilder,
	}
}

// CreateJob - Create a new job specification
func (s *JobsApiService) CreateJob(ctx context.Context, user string, jobSpec openapi.JobSpec) (openapi.ImplResponse, error) {
	jobStatus, err := s.dbService.CreateJob(user, jobSpec)
	if err != nil {
		errMsg := fmt.Errorf("failed to create a new job: %v", err)
		return errMsgFunc(errMsg)
	}

	dirty := false
	rollbackFunc := func() {
		err1 := s.dbService.DeleteTasks(jobStatus.Id, dirty)
		err2 := s.dbService.DeleteJob(user, jobStatus.Id)

		zap.S().Infof("delete tasks's error: %v; delete job's error: %v", err1, err2)
	}

	err = s.createTasks(user, jobStatus.Id, dirty)
	if err != nil {
		rollbackFunc()
		return errMsgFunc(err)
	}

	return openapi.Response(http.StatusCreated, jobStatus), nil
}

// DeleteJob - Delete job specification
func (s *JobsApiService) DeleteJob(ctx context.Context, user string, jobId string) (openapi.ImplResponse, error) {
	// TODO - update DeleteJob with the required logic for this service method.
	// Add api_jobs_service.go to the .openapi-generator-ignore to avoid overwriting this service
	// implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("DeleteJob method not implemented")
}

// GetJob - Get a job specification
func (s *JobsApiService) GetJob(ctx context.Context, user string, jobId string) (openapi.ImplResponse, error) {
	// TODO - update GetJob with the required logic for this service method.
	// Add api_jobs_service.go to the .openapi-generator-ignore to avoid overwriting this service
	// implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, JobSpec{}) or use other options such as http.Ok ...
	//return Response(200, JobSpec{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("GetJob method not implemented")
}

// GetJobStatus - Get job status of a given jobId
func (s *JobsApiService) GetJobStatus(ctx context.Context, user string, jobId string) (openapi.ImplResponse, error) {
	// TODO - update GetJobStatus with the required logic for this service method.
	// Add api_jobs_service.go to the .openapi-generator-ignore to avoid overwriting this service
	// implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, JobStatus{}) or use other options such as http.Ok ...
	//return Response(200, JobStatus{}), nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return Response(0, Error{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("GetJobStatus method not implemented")
}

// GetJobs - Get status info on all the jobs owned by user
func (s *JobsApiService) GetJobs(ctx context.Context, user string, limit int32) (openapi.ImplResponse, error) {
	jobsStatus, err := s.dbService.GetJobs(user, limit)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get status of all jobs: %v", err)
		zap.S().Debug(errMsg)
		return openapi.Response(http.StatusInternalServerError, err), fmt.Errorf(errMsg)
	}

	return openapi.Response(http.StatusOK, jobsStatus), nil
}

// GetTask - Get a job task for a given job and agent
func (s *JobsApiService) GetTask(ctx context.Context, jobId string, agentId string, key string) (openapi.ImplResponse, error) {
	if key == "" {
		return openapi.Response(http.StatusBadRequest, nil), fmt.Errorf("key can't be empty")
	}

	taskMap, err := s.dbService.GetTask(jobId, agentId, key)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("failed to get task")
	}

	return openapi.Response(http.StatusOK, taskMap), nil
}

// GetTasksInfo - Get the info of tasks in a job
func (s *JobsApiService) GetTasksInfo(ctx context.Context, user string, jobId string, limit int32) (openapi.ImplResponse, error) {
	tasksInfo, err := s.dbService.GetTasksInfo(user, jobId, limit, false)
	if err != nil {
		errMsg := fmt.Sprintf("failed to get info of all tasks in a job %s: %v", jobId, err)
		zap.S().Debug(errMsg)
		return openapi.Response(http.StatusInternalServerError, err), fmt.Errorf(errMsg)
	}

	return openapi.Response(http.StatusOK, tasksInfo), nil
}

// UpdateJob - Update a job specification
func (s *JobsApiService) UpdateJob(ctx context.Context, user string, jobId string,
	jobSpec openapi.JobSpec) (openapi.ImplResponse, error) {
	// Note: UpdateJob is a best-effort function
	//       No rollback mechanism is supported, which means that a user should
	//       continuously (modify the job spec and) try UpdateJob call
	err := s.dbService.UpdateJob(user, jobId, jobSpec)
	if err != nil {
		errMsg := fmt.Errorf("failed to create a new job: %v", err)
		return errMsgFunc(errMsg)
	}

	dirty := true
	err = s.createTasks(user, jobId, dirty)
	if err != nil {
		return errMsgFunc(err)
	}

	// delete unmodified tasks
	dirty = false
	err = s.dbService.DeleteTasks(jobId, dirty)
	if err != nil {
		return errMsgFunc(err)
	}

	// set dirty flag to false
	err = s.dbService.SetTaskDirtyFlag(jobId, dirty)
	if err != nil {
		return errMsgFunc(err)
	}

	return openapi.Response(http.StatusOK, nil), nil
}

// UpdateJobStatus - Update the status of a job
func (s *JobsApiService) UpdateJobStatus(ctx context.Context, user string, jobId string,
	jobStatus openapi.JobStatus) (openapi.ImplResponse, error) {
	// override jobId in the jobStatus
	jobStatus.Id = jobId

	event := job.NewJobEvent(user, jobStatus)
	s.jobEventQ.Enqueue(event)

	select {
	case <-time.After(defaultWaitTime):
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("response timed out")

	case err := <-event.ErrCh:
		if err != nil {
			errMsg := fmt.Sprintf("failed to update job status to %s: %v", jobStatus.State, err)
			zap.S().Debug(errMsg)
			return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf(errMsg)
		}

		return openapi.Response(http.StatusOK, nil), nil
	}
}

// UpdateTaskStatus - Update the status of a task
func (s *JobsApiService) UpdateTaskStatus(ctx context.Context, jobId string, agentId string,
	taskStatus openapi.TaskStatus) (openapi.ImplResponse, error) {
	err := s.dbService.UpdateTaskStatus(jobId, agentId, taskStatus)
	if err != nil {
		errMsg := fmt.Sprintf("failed to update a task status: %v", err)
		zap.S().Debug(errMsg)
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf(errMsg)
	}

	return openapi.Response(http.StatusCreated, nil), nil
}

func (s *JobsApiService) createTasks(user string, jobId string, dirty bool) error {
	// Obtain job specification
	jobSpec, err := s.dbService.GetJob(user, jobId)
	if err != nil {
		return fmt.Errorf("failed to get a job spec for job %s: %v", jobId, err)
	}

	tasks, _, err := s.jobBuilder.GetTasks(&jobSpec)
	if err != nil {
		return fmt.Errorf("failed to generate tasks for job %s: %v", jobId, err)
	}

	err = s.dbService.CreateTasks(tasks, dirty)
	if err != nil {
		return fmt.Errorf("failed to create tasks for job %s in database: %v", jobId, err)
	}

	return nil
}

func errMsgFunc(err error) (openapi.ImplResponse, error) {
	zap.S().Debug(err)
	return openapi.Response(http.StatusInternalServerError, nil), err
}
