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
		        -- SUM(cost) AS cost,
		        IFNULL(SUM(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END),0) AS cost_at_list,
		        IFNULL(SUM(cost),0) AS cost,
		      --  IFNULL(SUM(cost - IFNULL(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END, 0)),0) AS Negotiated_Savings
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@savings_programs_types))),0) AS Savings_Programs,
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@other_savings_types))),0) AS Other_Savings
		    FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		    WHERE
		        invoice.month BETWEEN @start_date AND @end_date
		        AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		        AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		        AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		        AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
				AND (ARRAY_LENGTH(@keys) IS NULL OR ARRAY_LENGTH(@keys) = 0 OR labels[SAFE_OFFSET(0)].key IN UNNEST(@keys))
                AND (ARRAY_LENGTH(@value) IS NULL OR ARRAY_LENGTH(@value) = 0 OR labels[SAFE_OFFSET(0)].value IN UNNEST(@value))
                AND (ARRAY_LENGTH(@region) IS NULL OR ARRAY_LENGTH(@region) = 0 OR location.region IN UNNEST(@region))
		    GROUP BY sku_id
		),
		cost_summary_last AS (
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
		        -- SUM(cost) AS cost,
		        IFNULL(SUM(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END),0) AS cost_at_list,
		        IFNULL(SUM(cost),0) AS cost,
		      --  IFNULL(SUM(cost - IFNULL(CASE WHEN sku.id = 'A9A8-879F-CB2C' THEN cost ELSE cost_at_list END, 0)),0) AS Negotiated_Savings
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@savings_programs_types))),0) AS Savings_Programs,
		        IFNULL(SUM((SELECT SUM(c.amount) FROM UNNEST(credits) AS c WHERE c.type IN UNNEST(@other_savings_types))),0) AS Other_Savings
		    FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
		    WHERE
		        invoice.month BETWEEN @prev_start AND @prev_end
		        AND _PARTITIONTIME BETWEEN TIMESTAMP(@prev_PartitionStartTime) AND TIMESTAMP(@prev_PartitionEndTime)
		        AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		        AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		        AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
				AND (ARRAY_LENGTH(@keys) IS NULL OR ARRAY_LENGTH(@keys) = 0 OR labels[SAFE_OFFSET(0)].key IN UNNEST(@keys))
                AND (ARRAY_LENGTH(@value) IS NULL OR ARRAY_LENGTH(@value) = 0 OR labels[SAFE_OFFSET(0)].value IN UNNEST(@value))
                AND (ARRAY_LENGTH(@region) IS NULL OR ARRAY_LENGTH(@region) = 0 OR location.region IN UNNEST(@region))
		    GROUP BY sku_id
		)
		SELECT
		    a.service_description,
		    a.sku_description,
		    a.sku_id,
		    FORMAT('%s %s',CAST(a.usage_amount AS STRING),CAST(a.usage_pricing_unit AS STRING)) AS Usage,
		    FORMAT('%s %s',CAST(a.usage_amount_details AS STRING),CAST(a.usage_pricing_unit_details AS STRING)) AS Usage_details,
		    -- 总使用成本
		    a.cost_at_list AS Usage_Cost,
		    -- Negotiated Savings
		    a.cost - a.cost_at_list AS Negotiated_Savings,
		    -- 各类折扣
		    a.Savings_Programs,
		    a.Other_Savings,
		    -- Sub_Total
		    a.cost_at_list + a.cost - a.cost_at_list + a.Savings_Programs + a.Other_Savings AS Sub_Total,
		    --  b.cost_at_list + b.cost - b.cost_at_list + b.Savings_Programs + b.Other_Savings AS Sub_Total,
		    CONCAT(
		      FORMAT(
		        '%.0f',
		        CASE
		          WHEN ROUND(b.cost_at_list + b.cost - b.cost_at_list + b.Savings_Programs + b.Other_Savings, 2) = 0
		          THEN 0  -- 分母为 0
		          ELSE
		            (
		              ROUND(a.cost_at_list + a.cost - a.cost_at_list + a.Savings_Programs + a.Other_Savings, 2)
		              -
		              ROUND(b.cost_at_list + b.cost - b.cost_at_list + b.Savings_Programs + b.Other_Savings, 2)
		            )
		            /
		            ROUND(b.cost_at_list + b.cost - b.cost_at_list + b.Savings_Programs + b.Other_Savings, 2)
		            * 100
		        END
		      ),
		      '%'
		    ) AS change_rate
		FROM cost_summary a
		LEFT JOIN cost_summary_last b USING (sku_id)
		ORDER BY sku_id;

			`

	return sql
}
