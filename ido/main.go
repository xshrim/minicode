package main

import "ido/service"

// OAuth规范: https://oauth.net/2/

// 注册应用
// http://gitlab.ebcpaas.com/profile/applications
// https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps

// API文档
// https://docs.gitlab.com/ce/api/issues.html
// https://docs.microsoft.com/en-us/graph/api/resources/outlooktask?view=graph-rest-beta

func main() {
	service.Server()
}
