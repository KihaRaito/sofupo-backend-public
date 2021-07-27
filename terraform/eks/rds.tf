resource "aws_db_parameter_group" "parameter_group" {
  name   = "sofupo-parameter-group"
  family = "postgres12"
}

resource "aws_subnet" "subnet1" {
  vpc_id     = module.vpc.vpc_id
  cidr_block = "10.0.4.0/24"

  tags = {
    Name = "sofupo-subnet1"
  }
}

resource "aws_subnet" "subnet2" {
  vpc_id     = module.vpc.vpc_id
  cidr_block = "10.0.5.0/24"

  tags = {
    Name = "sofupo-subnet2"
  }
}

resource "aws_subnet" "subnet3" {
  vpc_id     = module.vpc.vpc_id
  cidr_block = "10.0.6.0/24"

  tags = {
    Name = "sofupo-subnet3"
  }
}

resource "aws_db_subnet_group" "subnet_group" {
  name       = "sofupo-subnet-group"
  subnet_ids = [aws_subnet.subnet1.id, aws_subnet.subnet2.id, aws_subnet.subnet3.id]
}

resource "aws_db_instance" "db_instance" {
  identifier                 = "sofupo"
  engine                     = "postgres"
  engine_version             = "12.5"
  instance_class             = "db.t3.micro"
  allocated_storage          = 10
  name                       = "sofupo"
  username                   = var.db_username
  password                   = var.db_password
  multi_az                   = false
  publicly_accessible        = false
  auto_minor_version_upgrade = false
  deletion_protection        = false
  skip_final_snapshot        = true
  port                       = var.db_port
  apply_immediately          = false
  vpc_security_group_ids     = [aws_security_group.postgres_sg.id]
  parameter_group_name       = aws_db_parameter_group.parameter_group.name
  db_subnet_group_name       = aws_db_subnet_group.subnet_group.name

  lifecycle {
    ignore_changes = [password]
  }
}

resource "aws_security_group" "postgres_sg" {
    name = "postgres-sg"
    vpc_id = module.vpc.vpc_id
    ingress {
        from_port = var.db_port
        to_port = var.db_port
        protocol = "tcp"
        cidr_blocks = [module.vpc.vpc_cidr_block]
    }
    egress {
        from_port = 0
        to_port = 0
        protocol = "-1"
        cidr_blocks = ["0.0.0.0/0"]
    }
}

output "rds_endpoint" {
    value = aws_db_instance.db_instance.address
}
