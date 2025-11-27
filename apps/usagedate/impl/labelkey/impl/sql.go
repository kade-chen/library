package impl

// 6.83 GB/0 sec
func (s *service) queryByDateProjectLabelKeySQL() (sql string) {
	sql = `
		SELECT
		  project.id AS project_id,
		  labels[SAFE_OFFSET(0)].key as keys,
		  labels[SAFE_OFFSET(0)].value as value,
		  CONCAT(ANY_VALUE(labels[SAFE_OFFSET(0)].key),'/',ANY_VALUE(labels[SAFE_OFFSET(0)].value)) AS key_value_path
		  -- ANY_VALUE(service.description) AS service_description,
		  -- CONCAT('services/', service.id) AS service_path,
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		WHERE
		  DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		  AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		  AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
			  -- ⭐ 过滤 NULL 和 空字符串
		  AND labels[SAFE_OFFSET(0)].key IS NOT NULL
		  AND labels[SAFE_OFFSET(0)].value IS NOT NULL
		  AND labels[SAFE_OFFSET(0)].key <> ""
		  AND labels[SAFE_OFFSET(0)].value <> ""
		GROUP BY
		  project_id,keys,value
			`
	return sql
}
