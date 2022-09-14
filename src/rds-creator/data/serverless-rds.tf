terraform {
  required_providers {
    alicloud = {
      source = "registry.terraform.io/aliyun/alicloud"
      #      version = "1.141.0"
    }
  }
}

provider "alicloud" {
  region = var.region
}


resource "alicloud_vpc" "example" {
  vpc_name   = "photoview_vpc"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "example" {
  vpc_id     = alicloud_vpc.example.id
  cidr_block = "172.16.0.0/21"
  zone_id    = var.zone_id
}

resource "alicloud_db_instance" "example" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_name            = var.instance_name
  instance_storage         = "30"
  instance_type            = "mysql.n2.serverless.1c"
  vswitch_id               = alicloud_vswitch.example.id
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "SERVERLESS"
  monitoring_period        = "60"
  zone_id                  = var.zone_id
  whitelist_network_type   = "VPC"
  security_ips             = ["0.0.0.0/0",]
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
  #  instance_id = var.db_id
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

variable "zone_id" {
  description = "zone id"
  type        = string
  default     = "cn-hangzhou-i"
}


variable "region" {
  description = "resource region"
  type        = string
  default     = "cn-hangzhou"
}

output "DB_ID" {
  value = alicloud_db_instance.example.id
}
output "RESOURCE_ID" {
  value = alicloud_db_instance.example.resource_group_id
}

output "VPC_ID" {
  value       = alicloud_vpc.example.id
  description = "vpc id"
}


output "VSWITCH_ID" {
  value       = alicloud_vswitch.example.id
  description = "VSwitch Id"
}

output "SECURITY_GROUP_ID" {
  value       = alicloud_security_group.group.id
  description = "security_group"
}
