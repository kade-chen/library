package impl

func (s *service) queryByDateProjectSUSQL(projectIDs []string) (sql string) {
	if len(projectIDs) == 0 {
		// 查询全部项目
		sql = `
			SELECT
				service.id AS service_id,
				service.description AS service_description,
				CONCAT('services/', service.id) AS service_path,
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE
				DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) >= @start_date
				AND DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) <= @end_date
			GROUP BY
				service_id,
				service_description
			`
	} else {
		// 查询指定项目
		sql = `
			SELECT
				service.id AS service_id,
				service.description AS service_description,
				CONCAT('services/', service.id) AS service_path,
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) >= @start_date
				AND DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) <= @end_date
				AND project.id IN UNNEST(@project_ids)
			GROUP BY
				service_id,
				service_description;
			`
	}
	return sql
}
