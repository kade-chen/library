package impl

// 查询全部或者指定项目基于日期的project 6.12.17 GB / 1 sec / 126121 Slot milliseconds
func (s *service) queryBySkuSQL() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		  SELECT
		      -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		    ANY_VALUE(service.description) AS service_description, 
		    ANY_VALUE(sku.description) AS sku_description, 
		    sku.id AS sku_id, 
		    FORMAT('%.2f', IFNULL(SUM(usage.amount_in_pricing_units), 0)) AS usage_amount,
		    --   CAST(ROUND(IFNULL(SUM(usage.amount_in_pricing_units), 0), 10) AS STRING) AS usage_amount,
		    ANY_VALUE(usage.pricing_unit) AS usage_pricing_unit,
		    FORMAT('%.0f', IFNULL(SUM(usage.amount), 0)) AS usage_amount_details,
		    -- CAST(ROUND(IFNULL(SUM(usage.amount), 0), 10) AS STRING) AS usage_amount_details,
		    ANY_VALUE(usage.unit) AS usage_pricing_unit_details,
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
		  GROUP BY sku_id
		),

		-- 聚合指定 inovice service 的 cost
		service_summary AS (
		  SELECT
		    -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		    service.description AS service_description,
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
		  GROUP BY service_description
		),

		Saving_summary AS (
		  SELECT
		  -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		  sku.id AS sku_id, 
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
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),

		other_summary AS (
		  SELECT
		  -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		  sku.id AS sku_id, 
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
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),
		--------------------
		last_cost_summary AS (
		  SELECT
		        sku.id AS sku_id, 
		    SUM(cost) AS cost,
			  SUM(cost_at_list) AS cost_at_list
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		  GROUP BY sku_id
		),

		-- 聚合指定 inovice service 的 cost
		last_service_summary AS (
		  SELECT
		    service.description AS service_description,
			  ABS(SUM(cost_at_list)) AS invoice_cost_at_list_abs,
			  SUM(cost) AS invoice_cost
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		      AND service.id = "A656-35D2-EF7F"
		  GROUP BY service_description
		),

		last_Saving_summary AS (
		  SELECT
		  sku.id AS sku_id, 
		  SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
		  UNNEST(credits) AS credit
		ON TRUE
		WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),

		last_other_summary AS (
		  SELECT
		     sku.id AS sku_id, 
		  SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
		  UNNEST(credits) AS credit
		ON TRUE
		WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		)


		-- 合并结果
		SELECT
		    a.sku_description,
		    a.service_description,
		    a.sku_id,
		    FORMAT('%s %s',CAST(a.usage_amount AS STRING),CAST(a.usage_pricing_unit AS STRING)) AS Usage,
		    FORMAT('%s %s',CAST(a.usage_amount_details AS STRING),CAST(a.usage_pricing_unit_details AS STRING)) AS Usage_details,
			  -- a.usage_date,
			--   a.cost_at_list,
			--   b.invoice_cost_at_list_abs,
			--   b.invoice_cost,
			IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) AS Usage_Cost,
			
			CASE 
			    WHEN IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0) = 0 
			    THEN 0 
			    ELSE -ABS(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0)) 
			END AS Negotiated_Savings,
			
			CASE 
			    WHEN IFNULL(c.Savings_Programs, 0) = 0 
			    THEN 0 
			    ELSE -ABS(IFNULL(c.Savings_Programs, 0)) 
			END AS Savings_Programs,
			
			CASE 
			    WHEN IFNULL(d.Other_Savings, 0) = 0 
			    THEN 0 
			    ELSE -ABS(IFNULL(d.Other_Savings, 0)) 
			END AS Other_Savings,
			
			(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0))
			- ABS(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0))
			- ABS(IFNULL(c.Savings_Programs, 0))
			- ABS(IFNULL(d.Other_Savings, 0)) AS Sub_Total,
		    IFNULL(
		      FORMAT(
		        '%d%%',
		        CAST(
		          ROUND(
		            (
		              ROUND(
		                (IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0))
		                - ABS(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0))
		                - ABS(IFNULL(c.Savings_Programs, 0))
		                - ABS(IFNULL(d.Other_Savings, 0)),
		                2
		              )
		              -
		              ROUND(
		                (IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0))
		                - ABS(IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0) - IFNULL(aa.cost, 0))
		                - ABS(IFNULL(cc.Savings_Programs, 0))
		                - ABS(IFNULL(dd.Other_Savings, 0)),
		                2
		              )
		            )
		            / NULLIF(
		              ROUND(
		                (IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0))
		                - ABS(IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0) - IFNULL(aa.cost, 0))
		                - ABS(IFNULL(cc.Savings_Programs, 0))
		                - ABS(IFNULL(dd.Other_Savings, 0)),
		                2
		              ),
		              0
		            ) * 100  -- 转为百分比
		          ) AS INT64
		        )
		      ),
		      '0%'
		    ) AS change_rate
		FROM cost_summary a
		LEFT JOIN service_summary b USING (service_description)
		LEFT JOIN Saving_summary c USING (sku_id)
		LEFT JOIN other_summary d USING (sku_id)
		LEFT JOIN last_cost_summary aa USING (sku_id)
		LEFT JOIN last_service_summary bb USING (service_description)
		LEFT JOIN last_Saving_summary cc USING (sku_id)
		LEFT JOIN last_other_summary dd USING (sku_id)
		ORDER BY a.sku_id;
			`

	return sql
}


func (s *service) queryBySkuSQLBackup() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		  SELECT
		      -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		    ANY_VALUE(service.description) AS service_description, 
		    ANY_VALUE(sku.description) AS sku_description, 
		    sku.id AS sku_id, 
		    FORMAT('%.2f', IFNULL(SUM(usage.amount_in_pricing_units), 0)) AS usage_amount,
		    --   CAST(ROUND(IFNULL(SUM(usage.amount_in_pricing_units), 0), 10) AS STRING) AS usage_amount,
		    ANY_VALUE(usage.pricing_unit) AS usage_pricing_unit,
		    FORMAT('%.0f', IFNULL(SUM(usage.amount), 0)) AS usage_amount_details,
		    -- CAST(ROUND(IFNULL(SUM(usage.amount), 0), 10) AS STRING) AS usage_amount_details,
		    ANY_VALUE(usage.unit) AS usage_pricing_unit_details,
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
		  GROUP BY sku_id
		),

		-- 聚合指定 inovice service 的 cost
		service_summary AS (
		  SELECT
		    -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		    service.description AS service_description,
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
		  GROUP BY service_description
		),

		Saving_summary AS (
		  SELECT
		  -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		  sku.id AS sku_id, 
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
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),

		other_summary AS (
		  SELECT
		  -- DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
		  sku.id AS sku_id, 
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
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),
		--------------------
		last_cost_summary AS (
		  SELECT
		        sku.id AS sku_id, 
		    SUM(cost) AS cost,
			  SUM(cost_at_list) AS cost_at_list
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		  GROUP BY sku_id
		),

		-- 聚合指定 inovice service 的 cost
		last_service_summary AS (
		  SELECT
		    service.description AS service_description,
			  ABS(SUM(cost_at_list)) AS invoice_cost_at_list_abs,
			  SUM(cost) AS invoice_cost
		  FROM 
		      vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		  WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		      AND service.id = "A656-35D2-EF7F"
		  GROUP BY service_description
		),

		last_Saving_summary AS (
		  SELECT
		  sku.id AS sku_id, 
		  SUM(IF(credit.type IN UNNEST(@savings_programs_types), credit.amount, 0)) AS Savings_Programs
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
		  UNNEST(credits) AS credit
		ON TRUE
		WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		),

		last_other_summary AS (
		  SELECT
		     sku.id AS sku_id, 
		  SUM(IF(credit.type IN UNNEST(@other_savings_types), credit.amount, 0)) AS Other_Savings
		FROM
		  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		LEFT JOIN
		  UNNEST(credits) AS credit
		ON TRUE
		WHERE
		      DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @prev_start AND @prev_end
		      AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		      AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		      AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		      AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
		GROUP BY  sku_id
		-- ORDER BY usage_date DESC
		)


		-- 合并结果
		SELECT
		    a.sku_description,
		    a.service_description,
		    a.sku_id,
		    FORMAT('%s %s',CAST(a.usage_amount AS STRING),CAST(a.usage_pricing_unit AS STRING)) AS Usage,
		    FORMAT('%s %s',CAST(a.usage_amount_details AS STRING),CAST(a.usage_pricing_unit_details AS STRING)) AS Usage_details,
			  -- a.usage_date,
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
		      2) AS Sub_Total,
		    IFNULL(
		      FORMAT(
		        '%d%%',
		        CAST(
		          ROUND(
		            (
		              ROUND(
		                (IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0))
		                - ABS(IFNULL(a.cost_at_list, 0) + IFNULL(b.invoice_cost_at_list_abs, 0) + IFNULL(b.invoice_cost, 0) - IFNULL(a.cost, 0))
		                - ABS(IFNULL(c.Savings_Programs, 0))
		                - ABS(IFNULL(d.Other_Savings, 0)),
		                2
		              )
		              -
		              ROUND(
		                (IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0))
		                - ABS(IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0) - IFNULL(aa.cost, 0))
		                - ABS(IFNULL(cc.Savings_Programs, 0))
		                - ABS(IFNULL(dd.Other_Savings, 0)),
		                2
		              )
		            )
		            / NULLIF(
		              ROUND(
		                (IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0))
		                - ABS(IFNULL(aa.cost_at_list, 0) + IFNULL(bb.invoice_cost_at_list_abs, 0) + IFNULL(bb.invoice_cost, 0) - IFNULL(aa.cost, 0))
		                - ABS(IFNULL(cc.Savings_Programs, 0))
		                - ABS(IFNULL(dd.Other_Savings, 0)),
		                2
		              ),
		              0
		            ) * 100  -- 转为百分比
		          ) AS INT64
		        )
		      ),
		      '0%'
		    ) AS change_rate
		FROM cost_summary a
		LEFT JOIN service_summary b USING (service_description)
		LEFT JOIN Saving_summary c USING (sku_id)
		LEFT JOIN other_summary d USING (sku_id)
		LEFT JOIN last_cost_summary aa USING (sku_id)
		LEFT JOIN last_service_summary bb USING (service_description)
		LEFT JOIN last_Saving_summary cc USING (sku_id)
		LEFT JOIN last_other_summary dd USING (sku_id)
		ORDER BY a.sku_id;
			`

	return sql
}