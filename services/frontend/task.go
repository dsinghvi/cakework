package main

import (
	"database/sql"
	"time"

	"github.com/usecakework/cakework/lib/types"
)

func GetTaskLogs(db *sql.DB, userId string, appName string, taskName string, statusFilter string) (types.TaskLogs, error) {
	var requests []types.Request

	var rows *sql.Rows
	var err error

	if statusFilter == "" {
		rows, err = db.Query("SELECT requestId, status, parameters, result, createdAt, updatedAt FROM TaskRun where userId = ? AND app = ? AND task = ? ORDER BY updatedAt DESC LIMIT 100", userId, appName, taskName)
	} else {
		rows, err = db.Query("SELECT requestId, status, parameters, result, createdAt, updatedAt FROM TaskRun where userId = ? AND app = ? AND task = ? AND status = ? ORDER BY updatedAt DESC LIMIT 100", userId, appName, taskName, statusFilter)
	}

	if err != nil {
		return types.TaskLogs{
			Requests: requests,
		}, err
	}

	defer rows.Close()

	for rows.Next() {
		var result sql.NullString
		var createdAt time.Time
		var updatedAt time.Time
		var request types.Request
		if err := rows.Scan(&request.RequestId, &request.Status, &request.Parameters, &result, &createdAt, &updatedAt); err != nil {
			return types.TaskLogs{Requests: requests}, err
		}
		if result.Valid {
			request.Result = result.String
		}
		request.CreatedAt = createdAt.Unix()
		request.UpdatedAt = updatedAt.Unix()
		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return types.TaskLogs{
			Requests: requests,
		}, err
	}

	return types.TaskLogs{
		Requests: requests,
	}, nil
}
