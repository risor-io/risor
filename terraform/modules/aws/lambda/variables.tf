
variable "name" {
  description = "Function name"
  type        = string
  default     = ""
}

variable "handler" {
  description = "Function entrypoint"
  type        = string
  default     = ""
}

variable "package_path" {
  description = "Path to the function code"
  type        = string
}

variable "runtime" {
  description = "Function runtime"
  type        = string
  default     = ""
}

variable "role" {
  description = "IAM role to attach to the function"
  type        = string
  default     = ""
}

variable "description" {
  description = "Function description"
  type        = string
  default     = ""
}

variable "kms_key" {
  description = "KMS key to use for encryption operations"
  type        = string
  default     = null
}

variable "memory_size" {
  description = "Memory size in MB to assign to the function"
  type        = number
  default     = 128
}

variable "timeout" {
  description = "Function invocation timeout in seconds"
  type        = number
  default     = 3
}

variable "environment" {
  description = "Environment variables"
  type        = map(string)
  default     = {}
}

variable "tracing_mode" {
  description = "Function tracing mode"
  type        = string
  default     = null
}

variable "vpc_subnet_ids" {
  description = "List of subnet IDs to assign to the function"
  type        = list(string)
  default     = null
}

variable "vpc_security_group_ids" {
  description = "List of security group IDs to assign to the function"
  type        = list(string)
  default     = null
}

variable "tags" {
  description = "Tags to assign to the function"
  type        = map(string)
  default     = {}
}

variable "s3_bucket" {
  description = "Name of S3 bucket used to store the function code"
  type        = string
}

variable "s3_key" {
  description = "Path to the function code in the S3 bucket"
  type        = string
}
