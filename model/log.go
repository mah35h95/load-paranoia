package model

import "time"

type LoggingRequest struct {
	ProjectIDS    []string `json:"projectIds,omitempty"`
	ResourceNames []string `json:"resourceNames,omitempty"`
	Filter        string   `json:"filter,omitempty"`
	OrderBy       string   `json:"orderBy,omitempty"`
	PageSize      int32    `json:"pageSize,omitempty"`
	PageToken     string   `json:"pageToken,omitempty"`
}

type LoggingResponce struct {
	Entries       []Entry `json:"entries,omitempty"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
}

type Entry struct {
	ProtoPayload     ProtoPayload `json:"protoPayload,omitempty"`
	InsertID         string       `json:"insertId,omitempty"`
	Resource         Resource     `json:"resource,omitempty"`
	Timestamp        time.Time    `json:"timestamp,omitempty"`
	Severity         string       `json:"severity,omitempty"`
	LogName          string       `json:"logName,omitempty"`
	ReceiveTimestamp time.Time    `json:"receiveTimestamp,omitempty"`
}

type ProtoPayload struct {
	Type               string              `json:"@type,omitempty"`
	AuthenticationInfo AuthenticationInfo  `json:"authenticationInfo,omitempty"`
	RequestMetadata    RequestMetadata     `json:"requestMetadata,omitempty"`
	ServiceName        string              `json:"serviceName,omitempty"`
	MethodName         string              `json:"methodName,omitempty"`
	AuthorizationInfo  []AuthorizationInfo `json:"authorizationInfo,omitempty"`
	ResourceName       string              `json:"resourceName,omitempty"`
	ServiceData        ServiceData         `json:"serviceData,omitempty"`
}

type AuthenticationInfo struct {
	PrincipalEmail        string `json:"principalEmail,omitempty"`
	ServiceAccountKeyName string `json:"serviceAccountKeyName,omitempty"`
}

type AuthorizationInfo struct {
	Resource   string `json:"resource,omitempty"`
	Permission string `json:"permission,omitempty"`
	Granted    bool   `json:"granted,omitempty"`
}

type RequestMetadata struct {
	CallerIP                string `json:"callerIp,omitempty"`
	CallerSuppliedUserAgent string `json:"callerSuppliedUserAgent,omitempty"`
}

type ServiceData struct {
	Type                       string                     `json:"@type,omitempty"`
	JobGetQueryResultsResponse JobGetQueryResultsResponse `json:"jobGetQueryResultsResponse,omitempty"`
}

type JobGetQueryResultsResponse struct {
	Job Job `json:"job,omitempty"`
}

type Job struct {
	JobName          JobName          `json:"jobName,omitempty"`
	JobConfiguration JobConfiguration `json:"jobConfiguration,omitempty"`
	JobStatus        JobStatus        `json:"jobStatus,omitempty"`
	JobStatistics    JobStatistics    `json:"jobStatistics,omitempty"`
}

type JobConfiguration struct {
	Labels JobConfigurationLabels `json:"labels,omitempty"`
	Query  Query                  `json:"query,omitempty"`
}

type JobConfigurationLabels struct {
	DbtInvocationID string `json:"dbt_invocation_id,omitempty"`
}

type Query struct {
	Query             string `json:"query,omitempty"`
	DestinationTable  Table  `json:"destinationTable,omitempty"`
	CreateDisposition string `json:"createDisposition,omitempty"`
	WriteDisposition  string `json:"writeDisposition,omitempty"`
	QueryPriority     string `json:"queryPriority,omitempty"`
	StatementType     string `json:"statementType,omitempty"`
}

type Table struct {
	ProjectID string `json:"projectId,omitempty"`
	DatasetID string `json:"datasetId,omitempty"`
	TableID   string `json:"tableId,omitempty"`
}

type JobName struct {
	ProjectID string `json:"projectId,omitempty"`
	JobID     string `json:"jobId,omitempty"`
	Location  string `json:"location,omitempty"`
}

type JobStatistics struct {
	CreateTime           time.Time `json:"createTime,omitempty"`
	StartTime            time.Time `json:"startTime,omitempty"`
	EndTime              time.Time `json:"endTime,omitempty"`
	TotalProcessedBytes  string    `json:"totalProcessedBytes,omitempty"`
	TotalBilledBytes     string    `json:"totalBilledBytes,omitempty"`
	BillingTier          int64     `json:"billingTier,omitempty"`
	TotalSlotMS          string    `json:"totalSlotMs,omitempty"`
	ReferencedTables     []Table   `json:"referencedTables,omitempty"`
	TotalTablesProcessed int64     `json:"totalTablesProcessed,omitempty"`
	QueryOutputRowCount  string    `json:"queryOutputRowCount,omitempty"`
	Reservation          string    `json:"reservation,omitempty"`
}

type JobStatus struct {
	State string `json:"state,omitempty"`
}

type Resource struct {
	Type   string         `json:"type,omitempty"`
	Labels ResourceLabels `json:"labels,omitempty"`
}

type ResourceLabels struct {
	ProjectID string `json:"project_id,omitempty"`
}
