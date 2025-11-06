package impl

func (s *service) queryByDateProjectSKUsSQL(projectIDs []string) (sql string) {
	if len(projectIDs) == 0 {
		// 查询全部项目
		sql = `
			SELECT
				service.id AS service_id,
				sku.id AS sku_id,
				sku.description AS sku_describe,
				CONCAT('services/', service.id) AS service_path,
				CONCAT('services/', service.id, '/skus/', sku.id) AS service_sku_path,
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE
				DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) >= @start_date
				AND DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) <= @end_date
			GROUP BY
				service_id,
				sku_id,
				sku_describe
			`
	} else {
		// 查询指定项目
		sql = `
			SELECT
				service.id AS service_id,
				sku.id AS sku_id,
				sku.description AS sku_describe,
				CONCAT('services/', service.id) AS service_path,
				CONCAT('services/', service.id, '/skus/', sku.id) AS service_sku_path
			FROM vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) >= @start_date
				AND DATE(DATETIME(TIMESTAMP(_PARTITIONTIME), "America/Los_Angeles")) <= @end_date
				AND project.id IN UNNEST(@project_ids)
			GROUP BY
				service_id,
				sku_id,
				sku_describe;
			`
	}
	return sql
}
