package halon_smtpd_ctl

import (
	"fmt"
	"sort"
	"strings"
)

func (ps *ProcessStatsResponse) MarshalPrometheus() ([]byte, error) {
	var sb strings.Builder

	// Process runtime is a float, snowflake
	sb.WriteString(fmt.Sprintf("halon_process_runtime %.2f\n", ps.Process.Runtime))

	// Process
	// Resolver
	nameValues := map[string]uint64{
		"halon_process_pid":                   ps.Process.Pid,
		"halon_process_version_major":         ps.Process.Version.Major,
		"halon_process_version_minor":         ps.Process.Version.Minor,
		"halon_process_version_patch":         ps.Process.Version.Patch,
		"halon_resolver_pending":              ps.Resolver.Pending,
		"halon_resolver_dedup":                ps.Resolver.Dedup,
		"halon_resolver_running":              ps.Resolver.Running,
		"halon_resolver_maxrunning":           ps.Resolver.Maxrunning,
		"halon_resolver_cache_maxsize":        ps.Resolver.Cache.Maxsize,
		"halon_resolver_cache_size":           ps.Resolver.Cache.Size,
		"halon_resolver_cache_hits":           ps.Resolver.Cache.Hits,
		"halon_resolver_cache_misses":         ps.Resolver.Cache.Misses,
		"halon_resolver_cache_expires":        ps.Resolver.Cache.Expires,
		"halon_resolver_cache_evicts":         ps.Resolver.Cache.Evicts,
		"halon_resolver_cache_skips":          ps.Resolver.Cache.Skips,
		"halon_resolver_domain_cache_maxsize": ps.Resolver.Domain.Cache.Maxsize,
		"halon_resolver_domain_cache_size":    ps.Resolver.Domain.Cache.Size,
		"halon_resolver_domain_cache_hits":    ps.Resolver.Domain.Cache.Hits,
		"halon_resolver_domain_cache_misses":  ps.Resolver.Domain.Cache.Misses,
		"halon_resolver_domain_cache_expires": ps.Resolver.Domain.Cache.Expires,
		"halon_resolver_domain_cache_evicts":  ps.Resolver.Domain.Cache.Evicts,
	}

	// Servers
	for _, server := range ps.Servers {
		id := server.Serverid
		nameValues[fmt.Sprintf(`halon_servers_connections_concurrent{serverid="%s"}`, id)] = server.Connections.Concurrent
		nameValues[fmt.Sprintf(`halon_servers_connections_maxconcurrent{serverid="%s"}`, id)] = server.Connections.Maxconcurrent

		scripts := map[string]*ProcessStatsResponse_Script{
			"connect":    server.Scripts.Connect,
			"proxy":      server.Scripts.Proxy,
			"helo":       server.Scripts.Helo,
			"auth":       server.Scripts.Auth,
			"mailfrom":   server.Scripts.Mailfrom,
			"rcptto":     server.Scripts.Rcptto,
			"eod":        server.Scripts.Eod,
			"disconnect": server.Scripts.Disconnect,
		}
		for stage, script := range scripts {
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_pending{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Pending
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_running{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Running
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_finished{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Finished
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_errors{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Errors
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_threads{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Threads
			nameValues[fmt.Sprintf(`halon_servers_scripts_%s_suspended{serverid="%s",threadid="%s"}`, stage, id, script.Threadid)] = script.Suspended
		}
	}

	// Queue
	nameValues["halon_queue_loader_count"] = ps.Queue.Loader.Count
	nameValues["halon_queue_loader_pending"] = ps.Queue.Loader.Pending
	nameValues["halon_queue_loader_active"] = ps.Queue.Loader.Active
	nameValues["halon_queue_loader_maxactive"] = ps.Queue.Loader.Maxactive
	qs := map[string]*ProcessStatsResponse_Queue_Script{
		"predelivery":  ps.Queue.Scripts.Predelivery,
		"postdelivery": ps.Queue.Scripts.Postdelivery,
	}
	for stage, script := range qs {
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_pending{threadid="%s"}`, stage, script.Threadid)] = script.Pending
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_running{threadid="%s"}`, stage, script.Threadid)] = script.Running
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_finished{threadid="%s"}`, stage, script.Threadid)] = script.Finished
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_errors{threadid="%s"}`, stage, script.Threadid)] = script.Errors
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_threads{threadid="%s"}`, stage, script.Threadid)] = script.Threads
		nameValues[fmt.Sprintf(`halon_queue_scripts_%s_suspended{threadid="%s"}`, stage, script.Threadid)] = script.Suspended
	}
	nameValues["halon_queue_queue_defer_size"] = ps.Queue.Queue.Defer.Size
	nameValues["halon_queue_queue_active_size"] = ps.Queue.Queue.Active.Size
	nameValues["halon_queue_queue_active_size"] = ps.Queue.Queue.Active.Size
	for i, v := range ps.Queue.Queue.Active.Priorities {
		nameValues[fmt.Sprintf(`halon_queue_queue_active_priorities_size{priority="%d"}`, i)] = v.Size
	}
	nameValues["halon_queue_freeze_hold_size"] = ps.Queue.Freeze.Hold.Size
	nameValues["halon_queue_freeze_update_size"] = ps.Queue.Freeze.Update.Size
	nameValues["halon_queue_freeze_update_pending"] = ps.Queue.Freeze.Update.Pending
	nameValues["halon_queue_policy_concurrency_counters"] = ps.Queue.Policy.Concurrency.Counters
	nameValues["halon_queue_policy_concurrency_suspends"] = ps.Queue.Policy.Concurrency.Suspends
	nameValues["halon_queue_policy_rate_buckets"] = ps.Queue.Policy.Rate.Buckets
	nameValues["halon_queue_policy_rate_suspends"] = ps.Queue.Policy.Rate.Suspends
	nameValues["halon_queue_policy_dynamic_suspends"] = ps.Queue.Policy.Dynamic.Suspends
	nameValues["halon_queue_policy_dynamic_conditions"] = ps.Queue.Policy.Dynamic.Conditions
	nameValues["halon_queue_pickup_count"] = ps.Queue.Pickup.Count
	nameValues["halon_queue_pickup_skips"] = ps.Queue.Pickup.Skips
	nameValues["halon_queue_pickup_misses"] = ps.Queue.Pickup.Misses
	nameValues["halon_queue_pickup_pending"] = ps.Queue.Pickup.Pending
	nameValues["halon_queue_connections_concurrent"] = ps.Queue.Connections.Concurrent
	nameValues["halon_queue_connections_maxconcurrent"] = ps.Queue.Connections.Maxconcurrent
	nameValues["halon_queue_connections_pooling_size"] = ps.Queue.Connections.Pooling.Size
	nameValues["halon_queue_connections_pooling_maxsize"] = ps.Queue.Connections.Pooling.Maxsize
	nameValues["halon_queue_connections_pooling_hits"] = ps.Queue.Connections.Pooling.Hits
	nameValues["halon_queue_connections_pooling_misses"] = ps.Queue.Connections.Pooling.Misses
	nameValues["halon_queue_connections_pooling_expires"] = ps.Queue.Connections.Pooling.Expires
	nameValues["halon_queue_connections_pooling_evicts"] = ps.Queue.Connections.Pooling.Evicts
	nameValues["halon_queue_connections_pooling_skips"] = ps.Queue.Connections.Pooling.Skips
	nameValues["halon_queue_quota_size"] = ps.Queue.Quota.Size
	nameValues["halon_queue_delivery_delivered"] = ps.Queue.Delivery.Delivered
	nameValues["halon_queue_delivery_delayed"] = ps.Queue.Delivery.Delayed
	nameValues["halon_queue_delivery_failed"] = ps.Queue.Delivery.Failed
	nameValues["halon_queue_release_pending"] = ps.Queue.Release.Pending
	for _, script := range ps.Threads.Scripts {
		nameValues[fmt.Sprintf(`halon_threads_scripts_pending{id="%s"}`, script.Id)] = script.Pending
		nameValues[fmt.Sprintf(`halon_threads_scripts_rescheduled{id="%s"}`, script.Id)] = script.Rescheduled
		nameValues[fmt.Sprintf(`halon_threads_scripts_running{id="%s"}`, script.Id)] = script.Running
		nameValues[fmt.Sprintf(`halon_threads_scripts_maxrunning{id="%s"}`, script.Id)] = script.Maxrunning
		nameValues[fmt.Sprintf(`halon_threads_scripts_scripts{id="%s"}`, script.Id)] = script.Scripts
		nameValues[fmt.Sprintf(`halon_threads_scripts_maxscripts{id="%s"}`, script.Id)] = script.Maxscripts
	}

	var names []string
	for k, _ := range nameValues {
		names = append(names, k)
	}

	sort.Slice(names, func(i, j int) bool {
		return names[i] < names[j]
	})

	for _, k := range names {
		sb.WriteString(fmt.Sprintf("%s %d\n", k, nameValues[k]))
	}

	return []byte(sb.String()), nil
}
