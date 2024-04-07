terraform {
  required_providers {
    smtpd = {
      source = "smtpd.dev/smtpd/smtpd"
    }
  }
}

provider "smtpd" {
  api_key    = "foo@bar.com"
  api_secret = "foobar_password"
}


data "smtpd_profile" "test" {
  profile_name = "Daily Reports"
}
