
locals {
  azs = length(var.availability_zones) > 0 ? var.availability_zones : data.aws_availability_zones.available.names
}

data "aws_availability_zones" "available" {
  state = "available"
}

resource "aws_vpc" "this" {
  cidr_block           = var.cidr
  enable_dns_hostnames = var.enable_dns_hostnames
  enable_dns_support   = var.enable_dns_support
  tags                 = var.tags
}

resource "aws_subnet" "this" {
  count                   = length(var.subnets)
  cidr_block              = var.subnets[count.index]
  availability_zone       = local.azs[count.index]
  map_public_ip_on_launch = var.map_public_ip_on_launch
  vpc_id                  = aws_vpc.this.id
  tags                    = var.tags
}

resource "aws_route_table" "rt" {
  vpc_id = aws_vpc.this.id
  tags   = var.tags
}

resource "aws_route_table_association" "rt" {
  count          = length(var.subnets)
  subnet_id      = element(aws_subnet.this[*].id, count.index)
  route_table_id = aws_route_table.rt.id
}

resource "aws_route" "internet_gateway" {
  route_table_id         = aws_route_table.rt.id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.this.id
}

resource "aws_internet_gateway" "this" {
  vpc_id = aws_vpc.this.id
  tags   = var.tags
}
