package impl

// 查询全部或者指定项目基于日期的project 4.63 GB /1 sec / 88173 Slot milliseconds
func (s *service) queryByDateProjectSQL() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		    SELECT
		        DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		        ANY_VALUE(project.name) AS project_name,
		        project.id AS project_id,
		        ANY_VALUE(project.number) AS project_number,
		        -- cost_at_list: 对特定 SKU 使用 cost，其余使用 cost_at_list，直接处理 NULL
		        IFNULL(SUM(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END),0) AS cost_at_list,
		        IFNULL(SUM(cost),0) AS cost,
		      --  IFNULL(SUM(cost - IFNULL(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END, 0)),0) AS Negotiated_Savings
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@savings_programs_types))),0) AS Savings_Programs,
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@other_savings_types))),0) AS Other_Savings
		    FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		    WHERE
		        DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		        AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		        AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		        AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		        AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		    GROUP BY usage_date, project_id
		)
		SELECT
		    usage_date,
		    project_name,
		    project_id,
		    project_number, 
		    -- 总使用成本
		    cost_at_list AS Usage_Cost,
		    -- Negotiated Savings
		    cost - cost_at_list AS Negotiated_Savings,
		    -- 各类折扣
		    Savings_Programs,
		    Other_Savings,
		    -- Sub_Total
		    cost_at_list + cost - cost_at_list + Savings_Programs + Other_Savings AS Sub_Total
		FROM cost_summary
		ORDER BY usage_date DESC, project_id;
			`

	return sql
}

func (s *service) queryByDateProjectSQLBackup() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		  SELECT
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		      project.name AS project_name,
			    project.id AS project_id,
		      project.number AS project_number,
		    SUM(cost) AS cost,
			  SUM(cost_at_list) AS cost_at_list
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		  GROUP BY usage_date, project_id,project_number,project_name
		),

		-- 聚合指定 inovice service 的 cost
		service_summary AS (
		  SELECT
		    DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
			  project.id AS project_id,
			  ABS(SUM(cost_at_list)) AS invoice_cost_at_list_abs,
			  SUM(cost) AS invoice_cost
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		      AND service.id = "A656-35D2-EF7F"
		  GROUP BY usage_date, project_id
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
			--   a.cost_at_list,
			--   b.invoice_cost_at_list_abs,
			--   b.invoice_cost,
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
		      2) AS Sub_Total 
		FROM cost_summary a
		LEFT JOIN service_summary b USING (usage_date, project_id)
		LEFT JOIN Saving_summary c USING (usage_date, project_id)
		LEFT JOIN other_summary d USING (usage_date, project_id)
		ORDER BY a.usage_date DESC, a.project_id;
			`

	return sql
}