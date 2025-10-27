package impl

func (s *service) queryByDateProjectSQL(projectIDs []string) (sql string) {
	if len(projectIDs) == 0 {
		// 查询全部项目
		sql = `
			WITH cost_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				SUM(cost_at_list) AS cost_list
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) >= @start_date
				AND DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) <= @end_date
			GROUP BY usage_date, project_id
			),
			service_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				ABS(SUM(cost_at_list)) AS cost_list_abs,
				SUM(cost) AS cost
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) >= @start_date
			    AND DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) <= @end_date
				AND service.id = "A656-35D2-EF7F"
			GROUP BY usage_date, project_id
			)
			SELECT
				a.usage_date,
				a.project_id,
				a.cost_list,
				b.cost_list_abs,
				b.cost,
				IFNULL(a.cost_list,0) + IFNULL(b.cost_list_abs,0) + IFNULL(b.cost,0) AS total_cost_sum
			FROM cost_summary a
			LEFT JOIN service_summary b
			USING (usage_date, project_id)
			ORDER BY a.usage_date DESC, a.project_id;
			`
	} else {
		// 查询指定项目
		sql = `
			WITH cost_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				SUM(cost_at_list) AS cost_list
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) >= @start_date
			    AND DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) <= @end_date
				AND project.id IN UNNEST(@project_ids)
			GROUP BY usage_date, project_id
			),
			service_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				ABS(SUM(cost_at_list)) AS cost_list_abs,
				SUM(cost) AS cost
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) >= @start_date
			    AND DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) <= @end_date
				AND project.id IN UNNEST(@project_ids)
				AND service.id = "A656-35D2-EF7F"
			GROUP BY usage_date, project_id
			)
			SELECT
				a.usage_date,
				a.project_id,
				a.cost_list,
				b.cost_list_abs,
				b.cost,
				IFNULL(a.cost_list,0) + IFNULL(b.cost_list_abs,0) + IFNULL(b.cost,0) AS total_cost_sum
			FROM cost_summary a
			LEFT JOIN service_summary b
			USING (usage_date, project_id)
			ORDER BY a.usage_date DESC, a.project_id;
			`
	}
	return sql
}
