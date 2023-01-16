package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// TODO get frontend url from environment variable
// can share this with the poller

var (
	frontendURL string = "https://cakework-frontend.fly.dev"
)

type CreateTokenRequest struct {
	UserId string `json:"userId"`
	Name   string `json:"name"`
}

type CreateUserRequest struct {
	UserId string `json:"userId"`
}

type GetUserRequest struct {
	UserId string `json:"userId"`
}

type User struct {
	Id string `json:"id"`
}

type ClientToken struct {
	Token string `json:"token"`
}

// is there a way to avoid writing all this boilerplate?
func getUser(userId string, accessToken string, refreshToken string) *User { // TODO change return type
	url := frontendURL + "/get-user"
	getUserRequest := GetUserRequest{
		UserId: userId,
	}
	jsonReq, err := json.Marshal(getUserRequest)
	checkOsExit(err)

	req, err := newRequestWithAuth("GET", url, bytes.NewBuffer(jsonReq))
	checkOsExit(err)

	_, body, res := callHttp(req)
	if res.StatusCode == 200 {
		userId := body["id"].(string)
		return &User{Id: userId}
	} else {
		fmt.Println("Error getting user details")
		fmt.Println(res)
		return nil
	}
}

func createUser(userId string) *User { // TODO change return type
	url := frontendURL + "/create-user"
	getUserRequest := CreateUserRequest{
		UserId: userId,
	}
	jsonReq, err := json.Marshal(getUserRequest)
	checkOsExit(err)

	req, err := newRequestWithAuth("POST", url, bytes.NewBuffer(jsonReq))
	checkOsExit(err)

	_, body, res := callHttp(req)
	if res.StatusCode == 201 {
		userId := body["id"].(string)
		return &User{Id: userId}
	} else {
		fmt.Println("Error creating user")
		fmt.Println(res)
		return nil
	}
}

func createClientToken(userId string, name string) *ClientToken { // TODO change return type
	url := frontendURL + "/create-client-token"
	createTokenReq := CreateTokenRequest{
		UserId: userId,
		Name:   name,
	}
	jsonReq, err := json.Marshal(createTokenReq)
	checkOsExit(err)

	req, err := newRequestWithAuth("POST", url, bytes.NewBuffer(jsonReq))
	checkOsExit(err)

	_, body, res := callHttp(req)
	if res.StatusCode == 201 {
		token := body["token"].(string)
		return &ClientToken{Token: token}
	} else {
		fmt.Println("Error creating client token")
		fmt.Println(res)
		return nil
	}
}

type GetStatusRequest struct {
	UserId    string `json:"userId"`
	RequestId string `json:"requestId"`
}

// TODO return errors
func getRequestStatus(userId string, requestId string) string {
	url := frontendURL + "/get-status"
	getStatusRequest := GetStatusRequest{
		UserId:    userId,
		RequestId: requestId,
	}
	jsonReq, err := json.Marshal(getStatusRequest)
	checkOsExit(err)

	req, err := newRequestWithAuth("GET", url, bytes.NewBuffer(jsonReq))
	checkOsExit(err)

	_, body, res := callHttp(req)
	if res.StatusCode == 200 {
		status := body["status"].(string)
		return status
	} else if res.StatusCode == 404 {
		fmt.Println("Request ID " + requestId + " does not exist")
		return ""
	} else {
		checkOsExit(errors.New("Error getting request status, got an" + res.Status))
		return ""
	}
}

type GetTaskLogsRequest struct {
	UserId string `json:"userId"`
	App    string `json:"app"`
	Task   string `json:"task"`
}

type Request struct {
	RequestId  string `json:"request"`
	Status     string `json:"status"`
	Parameters string `json:"parameters"`
	Result     string `json:"result"`
}

type TaskLogs struct {
	Requests []Request `json:"requests"`
}

func getTaskLogs(userId string, appName string, taskName string, statuses []string) TaskLogs {
	url := frontendURL + "/task/logs"
	getTaskLogsRequest := GetTaskLogsRequest{
		UserId: userId,
		App:    appName,
		Task:   taskName,
	}
	jsonReq, err := json.Marshal(getTaskLogsRequest)
	checkOsExit(err)

	req, err := newRequestWithAuth("GET", url, bytes.NewBuffer(jsonReq))
	checkOsExit(err)

	res, err := http.DefaultClient.Do(req)
	checkOsExit(err)

	if res.StatusCode == 200 {
		var taskLogs TaskLogs
		bodybutbetter, err := io.ReadAll(res.Body)
		if err != nil {
			checkOsExit(errors.New("Error running task " + appName + "/" + taskName))
		}

		json.Unmarshal(bodybutbetter, &taskLogs)
		return taskLogs
	} else {
		// get res to string properly
		fmt.Println(res)
		checkOsExit(errors.New("Error running task " + appName + "/" + taskName))
		return TaskLogs{
			Requests: []Request{},
		}
	}
}
