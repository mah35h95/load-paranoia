package model

import "time"

const (
	OrderByAsc  string = "timestamp asc"
	OrderByDesc string = "timestamp desc"
)

type QueryLog struct {
	JobID          string
	OutputRowCount string
	From           time.Time
	To             time.Time
	TimestampFrom  time.Time
	TimestampTo    time.Time
	StartTime      time.Time
	EndTime        time.Time
}

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
	ProtoPayload     ProtoPayload `json:"protoPayload"`
	InsertID         string       `json:"insertId,omitempty"`
	Resource         Resource     `json:"resource"`
	Timestamp        time.Time    `json:"timestamp"`
	Severity         string       `json:"severity,omitempty"`
	LogName          string       `json:"logName,omitempty"`
	ReceiveTimestamp time.Time    `json:"receiveTimestamp"`
}

type ProtoPayload struct {
	Type               string              `json:"@type,omitempty"`
	AuthenticationInfo AuthenticationInfo  `json:"authenticationInfo"`
	RequestMetadata    RequestMetadata     `json:"requestMetadata"`
	ServiceName        string              `json:"serviceName,omitempty"`
	MethodName         string              `json:"methodName,omitempty"`
	AuthorizationInfo  []AuthorizationInfo `json:"authorizationInfo,omitempty"`
	ResourceName       string              `json:"resourceName,omitempty"`
	ServiceData        ServiceData         `json:"serviceData"`
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
	JobGetQueryResultsResponse JobGetQueryResultsResponse `json:"jobGetQueryResultsResponse"`
}

type JobGetQueryResultsResponse struct {
	Job Job `json:"job"`
}

type Job struct {
	JobName          JobName          `json:"jobName"`
	JobConfiguration JobConfiguration `json:"jobConfiguration"`
	JobStatus        JobStatus        `json:"jobStatus"`
	JobStatistics    JobStatistics    `json:"jobStatistics"`
}

type JobConfiguration struct {
	Labels JobConfigurationLabels `json:"labels"`
	Query  Query                  `json:"query"`
}

type JobConfigurationLabels struct {
	DbtInvocationID string `json:"dbt_invocation_id,omitempty"`
}

type Query struct {
	Query             string `json:"query,omitempty"`
	DestinationTable  Table  `json:"destinationTable"`
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
	CreateTime           time.Time `json:"createTime"`
	StartTime            time.Time `json:"startTime"`
	EndTime              time.Time `json:"endTime"`
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
	Labels ResourceLabels `json:"labels"`
}

type ResourceLabels struct {
	ProjectID string `json:"project_id,omitempty"`
}
