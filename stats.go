package redisethdb

// Stat typed string
type Stat string

const (
	/*
		Pool stats
	*/
	HITS         = "Hits"
	MISSES       = "Misses"
	TIMEOUTS     = "Timeouts"
	TOTAL_CONNS  = "TotalConns"
	IDLE_CONNS   = "IdleConns"
	STABLE_CONNS = "StaleConns"

	/*
		Data stats
	*/
	DB_SIZE = "DbSize"

	/*
		Info stats
	*/
	SERVER_INFO        = "Server"
	CLIENTS_INFO       = "Clients"
	MEMORY_INFO        = "Memory"
	PERSISTENCE_INFO   = "Persistence"
	STATS_INFO         = "Stats"
	REPLICATION_INFO   = "Replication"
	CPU_INFO           = "CPU"
	COMMAND_STATS_INFO = "Command"
	CLUSTER_INFO       = "Cluster"
	MODULES_INFO       = "Modules"
	KEYSPACE_INFO      = "Keyspace"
	ERROR_STATS_INFO   = "Errors"
	ALL_INFO           = "All"
	DEFAULT_INFO       = "Default"
	EVERYTHING         = "Everything"
)

func (s Stat) String() string {
	return string(s)
}

var poolStats = []Stat{
	HITS,
	MISSES,
	TIMEOUTS,
	TOTAL_CONNS,
	IDLE_CONNS,
	STABLE_CONNS,
}

var dataStats = []Stat{
	DB_SIZE,
}

var infoStats = []Stat{
	SERVER_INFO,
	CLIENTS_INFO,
	MEMORY_INFO,
	PERSISTENCE_INFO,
	STATS_INFO,
	REPLICATION_INFO,
	CPU_INFO,
	COMMAND_STATS_INFO,
	CLUSTER_INFO,
	MODULES_INFO,
	KEYSPACE_INFO,
	ERROR_STATS_INFO,
	ALL_INFO,
	DEFAULT_INFO,
	EVERYTHING,
}

func inList(str Stat, list []Stat) bool {
	for _, el := range list {
		if str == el {
			return true
		}
	}
	return false
}
