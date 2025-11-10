package apps

import (
	//impl config
	_ "github.com/kade-chen/google-billing-console/apps/configs/impl"

	//api config
	_ "github.com/kade-chen/google-billing-console/apps/project/api"


	//impl config
	_ "github.com/kade-chen/google-billing-console/apps/project/impl"
	_ "github.com/kade-chen/google-billing-console/apps/services/impl"
	_ "github.com/kade-chen/google-billing-console/apps/sku/impl"
)
