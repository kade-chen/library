package impl

import "fmt"

func (s *service) queryByDateProjectSKUsSQL1(table string) (sql string) {
	sql = `
			SELECT
			  service.id AS service_id,
			  sku.id AS sku_id,
			  ANY_VALUE(sku.description) AS sku_describe,
			  CONCAT('services/', service.id) AS service_path,
			  CONCAT('services/', service.id, '/skus/', sku.id) AS service_sku_path
			FROM
			  %v
			WHERE
			  invoice.month BETWEEN @start_date AND @end_date
			  AND usage_start_time BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
			  AND project.id IN UNNEST(@project_ids)
			GROUP BY
			  service_id,
			  sku_id
			`
	return fmt.Sprintf(sql, table)
}

// 6.83 GB/0 sec
func (s *service) queryByDateProjectSKUsSQL() (sql string) {
	sql = `
			SELECT
			  service.id AS service_id,
			  sku.id AS sku_id,
			  ANY_VALUE(sku.description) AS sku_describe,
			  CONCAT('services/', service.id) AS service_path,
			  CONCAT('services/', service.id, '/skus/', sku.id) AS service_sku_path
			FROM
			  vandorcloud-billing-check.detail_amount_all.gcp_billing_export_resource_v1_017C20_E02D28_86876D
			WHERE
			  invoice.month BETWEEN @start_date AND @end_date
			  AND _PARTITIONTIME BETWEEN TIMESTAMP(@PartitionStartTime) AND TIMESTAMP(@PartitionEndTime)
			  AND project.id IN UNNEST(@project_ids)
			GROUP BY
			  service_id,
			  sku_id
			`
	return sql
}
