package structs

import (
	"time"
)

const (
	QueryTTLMax = 24 * time.Hour
	QueryTTLMin = 10 * time.Second
)

// QueryDatacenterOptions sets options about how we fail over if there are no
// healthy nodes in the local datacenter.
type QueryDatacenterOptions struct {
	// NearestN is set to the number of remote datacenters to try, based on
	// network coordinates.
	NearestN int

	// Datacenters is a fixed list of datacenters to try after NearestN. We
	// never try a datacenter multiple times, so those are subtracted from
	// this list before proceeding.
	Datacenters []string
}

// QueryDNSOptions controls settings when query results are served over DNS.
type QueryDNSOptions struct {
	// TTL is the time to live for the served DNS results.
	TTL string
}

// ServiceQuery is used to query for a set of healthy nodes offering a specific
// service.
type ServiceQuery struct {
	// Service is the service to query.
	Service string

	// Failover controls what we do if there are no healthy nodes in the
	// local datacenter.
	Failover QueryDatacenterOptions

	// If OnlyPassing is true then we will only include nodes with passing
	// health checks (critical AND warning checks will cause a node to be
	// discarded)
	OnlyPassing bool

	// Tags are a set of required and/or disallowed tags. If a tag is in
	// this list it must be present. If the tag is preceded with "~" then
	// it is disallowed.
	Tags []string
}

// PreparedQuery defines a complete prepared query, and is the structure we
// maintain in the state store.
type PreparedQuery struct {
	// ID is this UUID-based ID for the query, always generated by Consul.
	ID string

	// Name is an optional friendly name for the query supplied by the
	// user. NOTE - if this feature is used then it will reduce the security
	// of any read ACL associated with this query/service since this name
	// can be used to locate nodes with supplying any ACL.
	Name string

	// Session is an optional session to tie this query's lifetime to. If
	// this is omitted then the query will not expire.
	Session string

	// Token is the ACL token used when the query was created, and it is
	// used when a query is subsequently executed. This token, or a token
	// with management privileges, must be used to change the query later.
	Token string

	// Service defines a service query (leaving things open for other types
	// later).
	Service ServiceQuery

	// DNS has options that control how the results of this query are
	// served over DNS.
	DNS QueryDNSOptions

	RaftIndex
}

type PreparedQueries []*PreparedQuery

type PreparedQueryOp string

const (
	PreparedQueryCreate PreparedQueryOp = "create"
	PreparedQueryUpdate                 = "update"
	PreparedQueryDelete                 = "delete"
)

// QueryRequest is used to create or change prepared queries.
type PreparedQueryRequest struct {
	Datacenter string
	Op         PreparedQueryOp
	Query      PreparedQuery
	WriteRequest
}

// RequestDatacenter returns the datacenter for a given request.
func (q *PreparedQueryRequest) RequestDatacenter() string {
	return q.Datacenter
}

// PreparedQueryExecuteRequest is used to execute a prepared query.
type PreparedQueryExecuteRequest struct {
	// Datacenter is the target this request is intended for.
	Datacenter string

	// QueryIDOrName is the ID of a query _or_ the name of one, either can
	// be provided.
	QueryIDOrName string

	// Limit will trim the resulting list down to the given limit.
	Limit int

	// Source is used to sort the results relative to a given node using
	// network coordinates.
	Source QuerySource

	// QueryOptions (unfortunately named here) controls the consistency
	// settings for the query lookup itself, as well as the service lookups.
	QueryOptions
}

// RequestDatacenter returns the datacenter for a given request.
func (q *PreparedQueryExecuteRequest) RequestDatacenter() string {
	return q.Datacenter
}

// PreparedQueryExecuteRemoteRequest is used when running a local query in a
// remote datacenter.
type PreparedQueryExecuteRemoteRequest struct {
	// Datacenter is the target this request is intended for.
	Datacenter string

	// Query is a copy of the query to execute.  We have to ship the entire
	// query over since it won't be present in the remote state store.
	Query PreparedQuery

	// Limit will trim the resulting list down to the given limit.
	Limit int

	// QueryOptions (unfortunately named here) controls the consistency
	// settings for the the service lookups.
	QueryOptions
}

// RequestDatacenter returns the datacenter for a given request.
func (q *PreparedQueryExecuteRemoteRequest) RequestDatacenter() string {
	return q.Datacenter
}

// PreparedQueryExecuteResponse has the results of executing a query.
type PreparedQueryExecuteResponse struct {
	Nodes CheckServiceNodes
	DNS   QueryDNSOptions
}
