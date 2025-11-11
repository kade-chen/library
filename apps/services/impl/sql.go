package impl

//3.51 GB/0 sec
func (s *service) queryByDateProjectSUSQL() (sql string) {
	sql = `
		SELECT
		  service.id AS service_id,
		  ANY_VALUE(service.description) AS service_description,
		  CONCAT('services/', service.id) AS service_path,
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		WHERE
		  DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		  AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		  AND project.id IN UNNEST(@project_ids)
		GROUP BY
		  service_id
			`
	return sql
}
