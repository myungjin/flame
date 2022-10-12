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

	"github.com/cisco-open/flame/cmd/controller/app/database"
	"github.com/cisco-open/flame/pkg/openapi"
)

// DesignSchemasApiService is a service that implents the logic for the DesignSchemasApiServicer
// This service should implement the business logic for every endpoint for the DesignSchemasApi API.
// Include any external packages or services that will be required by this service.
type DesignSchemasApiService struct {
	dbService database.DBService
}

// NewDesignSchemasApiService creates a default api service
func NewDesignSchemasApiService(dbService database.DBService) openapi.DesignSchemasApiServicer {
	return &DesignSchemasApiService{
		dbService: dbService,
	}
}

// CreateDesignSchema - Update a design schema
func (s *DesignSchemasApiService) CreateDesignSchema(ctx context.Context, user string, designId string,
	designSchema openapi.DesignSchema) (openapi.ImplResponse, error) {
	err := s.dbService.CreateDesignSchema(user, designId, designSchema)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("insert design schema details request failed")
	}

	return openapi.Response(http.StatusCreated, nil), nil
}

// GetDesignSchema - Get a design schema owned by user
func (s *DesignSchemasApiService) GetDesignSchema(ctx context.Context, user string, designId string,
	version string) (openapi.ImplResponse, error) {
	info, err := s.dbService.GetDesignSchema(user, designId, version)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design schema details request failed")
	}
	return openapi.Response(http.StatusOK, info), nil
}

// GetDesignSchemas - Get all design schemas in a design
func (s *DesignSchemasApiService) GetDesignSchemas(ctx context.Context, user string, designId string) (openapi.ImplResponse, error) {
	info, err := s.dbService.GetDesignSchemas(user, designId)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("get design schema details request failed")
	}
	return openapi.Response(http.StatusOK, info), nil
}

// UpdateDesignSchema - Update a schema for a given design
func (s *DesignSchemasApiService) UpdateDesignSchema(ctx context.Context, user string, designId string, version string,
	designSchema openapi.DesignSchema) (openapi.ImplResponse, error) {
	err := s.dbService.UpdateDesignSchema(user, designId, version, designSchema)
	if err != nil {
		return openapi.Response(http.StatusInternalServerError, nil), fmt.Errorf("schema update request failed")
	}

	return openapi.Response(http.StatusOK, nil), nil
}

func (s *DesignSchemasApiService) DeleteDesignSchema(ctx context.Context, user string, designId string,
	version string) (openapi.ImplResponse, error) {
	// TODO - update DeleteDesignCode with the required logic for this service method.
	// Add api_design_codes_service.go to the .openapi-generator-ignore
	// to avoid overwriting this service implementation when updating open api generation.

	//TODO: Uncomment the next line to return response Response(200, {}) or use other options such as http.Ok ...
	//return openapi.Response(200, nil),nil

	//TODO: Uncomment the next line to return response Response(404, {}) or use other options such as http.Ok ...
	//return openapi.Response(404, nil),nil

	//TODO: Uncomment the next line to return response Response(401, {}) or use other options such as http.Ok ...
	//return openapi.Response(401, nil),nil

	//TODO: Uncomment the next line to return response Response(0, Error{}) or use other options such as http.Ok ...
	//return openapi.Response(0, Error{}), nil

	return openapi.Response(http.StatusNotImplemented, nil), errors.New("DeleteDesignCode method not implemented")
}
