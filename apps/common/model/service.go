package model

import "cloud.google.com/go/bigquery"

type ServicesList struct {
	ServiceID   bigquery.NullString `bigquery:"service_id" json:"service_id"`
	ServiceDesc bigquery.NullString `bigquery:"service_description" json:"service_description"`
	ServicePath bigquery.NullString `bigquery:"service_path" json:"service_path"`
}
