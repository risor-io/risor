
variable "name" {
  description = "Function name"
  type        = string
  default     = ""
}

variable "description" {
  description = "Function description"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Tags to assign to the function"
  type        = map(string)
  default     = {}
}

variable "cidr" {
  description = "IPv4 CIDR block for the VPC"
  type        = string
  default     = "10.1.0.0/16"
}

variable "subnets" {
  description = "IPv4 CIDR blocks for the subnets"
  type        = list(string)
  default     = ["10.1.1.0/24"]
}

variable "availability_zones" {
  description = "Availability zones names in the region"
  type        = list(string)
  default     = []
}

variable "availability_zone_ids" {
  description = "Availability zone ids in the region"
  type        = list(string)
  default     = []
}

variable "enable_dns_hostnames" {
  description = "Enable DNS hostnames in the VPC"
  type        = bool
  default     = true
}

variable "enable_dns_support" {
  description = "Enable DNS support in the VPC"
  type        = bool
  default     = true
}

variable "map_public_ip_on_launch" {
  description = "Enable public IP addresses when launching instances"
  type        = bool
  default     = true
}
