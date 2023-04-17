provider "alicloud" {
  region = var.region
}


data "alicloud_db_zones" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "Serverless"
  category                 = "serverless_basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_fc_zones" "zones_ids" {}

data "alicloud_db_instance_classes" "example" {
  zone_id                  = data.alicloud_db_zones.example.ids.1
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "serverless_basic"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "Serverless"
  commodity_code           = "rds_serverless_public_cn"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "photoview_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id     = alicloud_vpc.example.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alicloud_db_zones.example.ids.1
}

resource "alicloud_vswitch" "fc_example" {
  vpc_id     = alicloud_vpc.example.id
  cidr_block = "172.16.3.0/24"
  zone_id    = data.alicloud_fc_zones.zones_ids.ids.1
}

resource "alicloud_db_instance" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_name            = var.instance_name
  monitoring_period        = "60"
  whitelist_network_type   = "VPC"
  security_ips             = ["0.0.0.0/0",]
  instance_storage         = data.alicloud_db_instance_classes.example.instance_classes.0.storage_range.min
  instance_type            = data.alicloud_db_instance_classes.example.instance_classes.0.instance_class
  instance_charge_type     = "Serverless"
  zone_id                  = data.alicloud_db_zones.example.ids.1
  vswitch_id               = alicloud_vswitch.example.id
  db_instance_storage_type = "cloud_essd"
  serverless_config {
    max_capacity = 8
    min_capacity = 0.5
    auto_pause   = false
    switch_force = false
  }
  category = "serverless_basic"
}


resource "alicloud_db_account" "example" {
  db_instance_id   = alicloud_db_instance.example.id
  #  db_instance_id = var.db_id
  account_name     = var.instance_name
  account_password = "Photoview2022"
  account_type     = "Super"
}


resource "alicloud_db_database" "default" {
  instance_id   = alicloud_db_instance.example.id
  name          = var.instance_name
  character_set = "utf8"
}


resource "alicloud_security_group" "group" {
  name   = var.instance_name
  vpc_id = alicloud_vpc.example.id
}


variable "instance_name" {
  description = "instance name"
  type = string
  default = "photoview"
}


variable "region" {
  description = "resource region"
  type        = string
  default     = "cn-hangzhou"
}

output "FC_VSWITCH_ID" {
  value = alicloud_vswitch.fc_example.id
}

output "DB_ID" {
  value = alicloud_db_instance.example.id
}
output "RESOURCE_ID" {
  value = alicloud_db_instance.example.resource_group_id
}

output "VPC_ID" {
  value       = alicloud_vpc.example.id
}


output "VSWITCH_ID" {
  value       = alicloud_vswitch.example.id
}

output "SECURITY_GROUP_ID" {
  value       = alicloud_security_group.group.id
}
