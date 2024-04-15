provider "aws" {
  region = "us-west-2" # Specify the AWS region
}

resource "aws_vpc" "example" {
  cidr_block = "10.0.0.0/16" # Specify the CIDR block for the VPC

  tags = {
    Name = "example-vpc"
  }
}