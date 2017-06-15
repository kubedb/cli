package util

import (
apiv1 "k8s.io/client-go/pkg/api/v1"
)

// ref: k8s.io/kubernetes/pkg/api/events/sorted_event_list.go

type SortableEvents []apiv1.Event

func (list SortableEvents) Len() int {
	return len(list)
}

func (list SortableEvents) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list SortableEvents) Less(i, j int) bool {
	return list[i].LastTimestamp.Time.After(list[j].LastTimestamp.Time)
}
