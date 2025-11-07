package impl

// 查询全部或者指定项目基于日期的project 6.19 GB/1 sec
func (s *service) queryByDateProjectSQL() (sql string) {
	// 查询全部项目
	sql = `
			WITH cost_summary AS (
				SELECT
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
					ANY_VALUE(project.name) AS project_name,
					project.id AS project_id,
					project.number AS project_number,
					SUM(cost) AS cost,
					SUM(cost_at_list) AS cost_at_list
				FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
				WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				GROUP BY usage_date, project_id, project_number
				),
			service_summary AS (
				SELECT
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
					project.id AS project_id,
					ABS(SUM(cost_at_list)) AS invoice_cost_at_list_abs,
					SUM(cost) AS invoice_cost
				FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
				WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
					AND service.id = "A656-35D2-EF7F"
				GROUP BY usage_date, project_id
				),
			Saving_summary AS (
				SELECT
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
					project.id AS project_id,
					SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
				FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
				LEFT JOIN
				UNNEST(credits) AS credit
				ON TRUE
				WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				GROUP BY usage_date, project_id
				ORDER BY usage_date DESC
				),
			other_summary AS (
				SELECT
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
					project.id AS project_id,
					SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
				FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
				LEFT JOIN
				UNNEST(credits) AS credit
				ON TRUE
				WHERE DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				GROUP BY usage_date, project_id
				ORDER BY usage_date DESC
				)
			SELECT
				a.usage_date,
				a.project_name,
				a.project_id,
				a.project_number,
				a.cost_at_list,
				b.invoice_cost_at_list_abs,
				b.invoice_cost,
				ROUND(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0), 2) AS Usage_Cost,

					CASE WHEN ROUND(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0), 2) = 0 
					THEN 0 
					ELSE -ABS(ROUND(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0), 2)) 
					END AS Negotiated_Savings,

					CASE WHEN ROUND(IFNULL(c.Savings_Programs, 0), 2) = 0 
						THEN 0 
						ELSE -ABS(ROUND(IFNULL(c.Savings_Programs, 0), 2)) 
					END AS Savings_Programs,

					CASE WHEN ROUND(IFNULL(d.Other_Savings, 0), 2) = 0 
						THEN 0 
						ELSE -ABS(ROUND(IFNULL(d.Other_Savings, 0), 2)) 
					END AS Other_Savings,
					
					ROUND(
					(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0))
					- ABS(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0))
					- ABS(IFNULL(c.Savings_Programs, 0))
					- ABS(IFNULL(d.Other_Savings, 0)),
					2) AS sub_total 
			FROM cost_summary a
			LEFT JOIN service_summary b USING (usage_date, project_id)
			LEFT JOIN Saving_summary c USING (usage_date, project_id)
			LEFT JOIN other_summary d USING (usage_date, project_id)
			ORDER BY a.usage_date DESC, a.project_id;
			`

	return sql
}

// 查询全部或者指定项目基于日期的project,全部service，指定skus， 4.46 GB/ 0 sec
func (s *service) queryByDateProjectServicesCustomSkusSQL() (sql string) {
	// 查询全部项目
	sql = `
			WITH cost_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				ANY_VALUE(project.name) AS project_name, 
				project.id AS project_id,
				project.number AS project_number,
				SUM(cost_at_list) AS Usage_Cost,

				CASE 
				WHEN IFNULL(SUM(cost_at_list) - SUM(cost), 0) = 0 
					THEN 0 
				ELSE -ABS(
					IFNULL(SUM(cost_at_list) - SUM(cost), 0)
				)
				END AS Negotiated_Savings,
				FROM 
					vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
				WHERE
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
					AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
				GROUP BY usage_date, project_id,project_number
			),

			Saving_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
			FROM
				vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			LEFT JOIN
				UNNEST(credits) AS credit
			ON TRUE
			WHERE
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
					AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
			GROUP BY usage_date, project_id
			ORDER BY usage_date DESC
			),

			other_summary AS (
			SELECT
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
				project.id AS project_id,
				SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
			FROM
				vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			LEFT JOIN
				UNNEST(credits) AS credit
			ON TRUE
			WHERE
					DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
					AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
					AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
					AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
			GROUP BY usage_date, project_id
			ORDER BY usage_date DESC
			)

			-- 合并结果
			SELECT
				a.usage_date,
				a.project_name,
				a.project_id,
				a.project_number,
				ROUND(a.Usage_Cost, 2) AS Usage_Cost,
				ROUND(a.Negotiated_Savings, 2) AS Negotiated_Savings,
				ROUND(b.Savings_Programs, 2) AS Savings_Programs,
				ROUND(c.Other_Savings, 2) AS Other_Savings, 
				ROUND(
				(
					(IFNULL(a.Usage_Cost, 0))
					- ABS(IFNULL(a.Negotiated_Savings, 0))
					- ABS(IFNULL(b.Savings_Programs, 0))
					- ABS(IFNULL(c.Other_Savings, 0))
				),2
				) AS Sub_Total
			FROM cost_summary a
			LEFT JOIN Saving_summary b USING (usage_date, project_id)
			LEFT JOIN other_summary c USING (usage_date, project_id)
			ORDER BY a.usage_date DESC, a.project_id;
			`
	return sql
}

// 查询全部或者指定项目基于日期的project,指定service，全sku，4.46 GB/0 sec
func (s *service) queryByDateProjectCustomServicesAllSkusSQL() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			ANY_VALUE(project.name) AS project_name, 
			project.id AS project_id,
			project.number AS project_number,
			SUM(cost_at_list) AS Usage_Cost,

			CASE 
			WHEN IFNULL(SUM(cost_at_list) - SUM(cost), 0) = 0 
				THEN 0 
			ELSE -ABS(
				IFNULL(SUM(cost_at_list) - SUM(cost), 0)
			)
			END AS Negotiated_Savings,
			FROM 
				vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
			GROUP BY usage_date, project_id,project_number
		),

		Saving_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			project.id AS project_id,
			SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
		FROM
			vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
			UNNEST(credits) AS credit
		ON TRUE
		WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		GROUP BY usage_date, project_id
		ORDER BY usage_date DESC
		),

		other_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			project.id AS project_id,
			SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
		FROM
			vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
			UNNEST(credits) AS credit
		ON TRUE
		WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		GROUP BY usage_date, project_id
		ORDER BY usage_date DESC
		)

		-- 合并结果
		SELECT
			a.usage_date,
			a.project_name,
			a.project_id,
			a.project_number,
			ROUND(a.Usage_Cost, 2) AS Usage_Cost,
			ROUND(a.Negotiated_Savings, 2) AS Negotiated_Savings,
			ROUND(b.Savings_Programs, 2) AS Savings_Programs,
			ROUND(c.Other_Savings, 2) AS Other_Savings, 
			ROUND(
			(
				(IFNULL(a.Usage_Cost, 0))
				- ABS(IFNULL(a.Negotiated_Savings, 0))
				- ABS(IFNULL(b.Savings_Programs, 0))
				- ABS(IFNULL(c.Other_Savings, 0))
			),2
			) AS Sub_Total
		FROM cost_summary a
		LEFT JOIN Saving_summary b USING (usage_date, project_id)
		LEFT JOIN other_summary c USING (usage_date, project_id)
		ORDER BY a.usage_date DESC, a.project_id;
			`
	return sql
}

// 查询全部或者指定项目基于日期的project,指定service，指定sku，5.49 GB /0 sec
func (s *service) queryByDateProjectCustomServicesCustomSkusSQL() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			ANY_VALUE(project.name) AS project_name, 
			project.id AS project_id,
			project.number AS project_number,
			SUM(cost_at_list) AS Usage_Cost,

			CASE 
			WHEN IFNULL(SUM(cost_at_list) - SUM(cost), 0) = 0 
				THEN 0 
			ELSE -ABS(
				IFNULL(SUM(cost_at_list) - SUM(cost), 0)
			)
			END AS Negotiated_Savings,
			FROM 
				vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
				AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
			GROUP BY usage_date, project_id,project_number
		),

		Saving_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			project.id AS project_id,
			SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
		FROM
			vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
			UNNEST(credits) AS credit
		ON TRUE
		WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
				AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY usage_date, project_id
		ORDER BY usage_date DESC
		),

		other_summary AS (
		SELECT
			DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			project.id AS project_id,
			SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
		FROM
			vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
			UNNEST(credits) AS credit
		ON TRUE
		WHERE
				DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
				AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
				AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
				AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
				AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY usage_date, project_id
		ORDER BY usage_date DESC
		)

		-- 合并结果
		SELECT
			a.usage_date,
			a.project_name,
			a.project_id,
			a.project_number,
			ROUND(a.Usage_Cost, 2) AS Usage_Cost,
			ROUND(a.Negotiated_Savings, 2) AS Negotiated_Savings,
			ROUND(b.Savings_Programs, 2) AS Savings_Programs,
			ROUND(c.Other_Savings, 2) AS Other_Savings, 
			ROUND(
			(
				(IFNULL(a.Usage_Cost, 0))
				- ABS(IFNULL(a.Negotiated_Savings, 0))
				- ABS(IFNULL(b.Savings_Programs, 0))
				- ABS(IFNULL(c.Other_Savings, 0))
			),2
			) AS Sub_Total
		FROM cost_summary a
		LEFT JOIN Saving_summary b USING (usage_date, project_id)
		LEFT JOIN other_summary c USING (usage_date, project_id)
		ORDER BY a.usage_date DESC, a.project_id;
			`

	return sql
}

func (s *service) queryByDateProjectServicesCustomSkusSQ1(projectIDs, serviceIDs []string) (sql string) {
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
				IFNULL(a.cost_list,0) + IFNULL(b.cost_list_abs,0) + IFNULL(b.cost,0) AS Usage_Cost
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
				IFNULL(a.cost_list,0) + IFNULL(b.cost_list_abs,0) + IFNULL(b.cost,0) AS Usage_Cost
			FROM cost_summary a
			LEFT JOIN service_summary b
			USING (usage_date, project_id)
			ORDER BY a.usage_date DESC, a.project_id;
			`
	}
	return sql
}
