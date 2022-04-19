/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package apis

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/chaosnative/chaosctl/pkg/utils"

	"github.com/chaosnative/chaosctl/pkg/types"
)

type createProjectResponse struct {
	Message string `json:"message"`
	Data    struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"data"`
	Errors []struct {
		Path []string `json:"path"`
	} `json:"errors"`
}

type createProjectPayload struct {
	ProjectName string `json:"project_name"`
	UserID      string `json:"user_id"`
}

func CreateProjectRequest(userID string, projectName string, cred types.Credentials) (createProjectResponse, error) {
	payloadBytes, err := json.Marshal(createProjectPayload{
		ProjectName: projectName,
		UserID:      userID,
	})

	if err != nil {
		return createProjectResponse{}, err
	}
	resp, err := SendRequest(SendRequestParams{cred.Endpoint + utils.AuthAPIPath + "/create_project", "Bearer " + cred.Token}, payloadBytes, string(types.Post))
	if err != nil {
		return createProjectResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return createProjectResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var project createProjectResponse
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return createProjectResponse{}, err
		}

		utils.White_B.Println("project/" + project.Data.Name + " created")
		return project, nil
	} else {
		return createProjectResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

type listProjectResponse struct {
	Data []struct {
		ID        string `json:"ID"`
		Name      string `json:"Name"`
		CreatedAt string `json:"CreatedAt"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

func ListProject(cred types.Credentials) (listProjectResponse, error) {

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.AuthAPIPath + "/list_projects", Token: "Bearer " + cred.Token}, []byte{}, string(types.Get))
	if err != nil {
		return listProjectResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return listProjectResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data listProjectResponse
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			return listProjectResponse{}, err
		}
		return data, nil
	} else {
		return listProjectResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

type ProjectDetails struct {
	Data Data `json:"data"`
}

type Data struct {
	ID         string    `json:"ID"`
	Username   string    `json:"Username"`
	CreatedAt  string    `json:"CreatedAt"`
	UserStatus string    `json:"UserStatus"`
	Email      string    `json:"Email"`
	FirstName  string    `json:"FirstName"`
	LastName   string    `json:"LastName"`
	Projects   []Project `json:"Projects"`
}

type Member struct {
	UserID     string `json:"UserID"`
	Role       string `json:"Role"`
	Invitation string `json:"Invitation"`
	JoinedAt   string `json:"JoinedAt"`
}

type Project struct {
	ID    string `json:"ID"`
	Name  string `json:"Name"`
	Owner struct {
		AccID   string `json:"AccID"`
		AccType string `json:"AccType"`
	} `json:"Owner"`
	Teams     []interface{} `json:"Teams"`
	State     interface{}   `json:"State"`
	CreatedAt string        `json:"CreatedAt"`
	UpdatedAt string        `json:"UpdatedAt"`
	RemovedAt string        `json:"RemovedAt"`
	Members   []Member      `json:"Members"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(c types.Credentials) (ProjectDetails, error) {
	token, _ := jwt.Parse(c.Token, nil)
	if token == nil {
		return ProjectDetails{}, nil
	}
	Username, _ := token.Claims.(jwt.MapClaims)["username"].(string)
	resp, err := SendRequest(SendRequestParams{Endpoint: c.Endpoint + utils.AuthAPIPath + "/get_user_with_project/" + Username, Token: "Bearer " + c.Token}, []byte{}, string(types.Get))
	if err != nil {
		return ProjectDetails{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return ProjectDetails{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var project ProjectDetails
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return ProjectDetails{}, err
		}

		return project, nil
	} else {
		return ProjectDetails{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
