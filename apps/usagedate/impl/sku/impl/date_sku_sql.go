package impl

// 查询全部或者指定项目基于日期的project 9.09 GB / 1 sec / 84374 Slot milliseconds
func (s *service) queryByDateSkuSQL() (sql string) {
	// 查询全部项目
	sql = `
		WITH cost_summary AS (
		    SELECT
		        DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) AS usage_date,
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
		        DATE(DATETIME(TIMESTAMP(usage_start_time), "America/Los_Angeles")) BETWEEN @start_date AND @end_date
		        AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
		        AND (ARRAY_LENGTH(@project_ids) IS NULL OR ARRAY_LENGTH(@project_ids) = 0 OR project.id IN UNNEST(@project_ids))
		        AND (ARRAY_LENGTH(@services_ids) IS NULL OR ARRAY_LENGTH(@services_ids) = 0 OR service.id IN UNNEST(@services_ids))
		        AND (ARRAY_LENGTH(@skus_ids) IS NULL OR ARRAY_LENGTH(@skus_ids) = 0 OR sku.id IN UNNEST(@skus_ids))
				AND (ARRAY_LENGTH(@keys) IS NULL OR ARRAY_LENGTH(@keys) = 0 OR labels[SAFE_OFFSET(0)].key IN UNNEST(@keys))
                AND (ARRAY_LENGTH(@value) IS NULL OR ARRAY_LENGTH(@value) = 0 OR labels[SAFE_OFFSET(0)].value IN UNNEST(@value))
                AND (ARRAY_LENGTH(@region) IS NULL OR ARRAY_LENGTH(@region) = 0 OR location.region IN UNNEST(@region))
		    GROUP BY usage_date, sku_id
		)
		SELECT
		    usage_date,
			sku_description,
		    service_description,
		    sku_id,
		    FORMAT('%s %s',CAST(usage_amount AS STRING),CAST(usage_pricing_unit AS STRING)) AS Usage,
		    FORMAT('%s %s',CAST(usage_amount_details AS STRING),CAST(usage_pricing_unit_details AS STRING)) AS Usage_details,
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
		ORDER BY usage_date DESC, sku_id;
			`

	return sql
}

