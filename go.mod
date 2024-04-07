module github.com/helloallan/terraform-provider-smtpd

go 1.16

require github.com/hashicorp/terraform-plugin-sdk/v2 v2.6.1 // indirect

replace (
	github.com/helloallan/terraform-provider-smtpd/smtpd => ./terraform-provider-smtpd/smtpd
	github.com/helloallan/terraform-provider-smtpd/sdk => ./terraform-provider-smtpd/sdk
)
