# Scheduled Disassembly scan on AWS Fargate, triggered by EventBridge.
# SARIF is written by the container; ship it to S3 via your report sink of choice.
resource "aws_ecs_cluster" "this" {
  name = var.name
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/ecs/${var.name}"
  retention_in_days = 30
}

data "aws_region" "current" {}

data "aws_iam_policy_document" "assume" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "exec" {
  name               = "${var.name}-exec"
  assume_role_policy = data.aws_iam_policy_document.assume.json
}

resource "aws_iam_role_policy_attachment" "exec" {
  role       = aws_iam_role.exec.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# Allow pulling the API key secret at task start.
resource "aws_iam_role_policy" "secret" {
  role = aws_iam_role.exec.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["secretsmanager:GetSecretValue"]
      Resource = var.api_key_secret_arn
    }]
  })
}

resource "aws_ecs_task_definition" "this" {
  family                   = var.name
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = "512"
  memory                   = "1024"
  execution_role_arn       = aws_iam_role.exec.arn
  container_definitions = jsonencode([{
    name        = "disassembly"
    image       = var.image
    command     = ["scan", var.target, "--ci"]
    secrets     = [{ name = "DISASSEMBLY_API_KEY", valueFrom = var.api_key_secret_arn }]
    environment = [{ name = "REPORTS_BUCKET", value = var.reports_bucket }]
    logConfiguration = {
      logDriver = "awslogs"
      options = {
        "awslogs-group"         = aws_cloudwatch_log_group.this.name
        "awslogs-region"        = data.aws_region.current.name
        "awslogs-stream-prefix" = "scan"
      }
    }
  }])
}

# EventBridge schedule -> RunTask
resource "aws_cloudwatch_event_rule" "schedule" {
  name                = "${var.name}-schedule"
  schedule_expression = var.schedule
}

resource "aws_iam_role" "events" {
  name = "${var.name}-events"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Allow"
      Principal = { Service = "events.amazonaws.com" }
      Action    = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "events" {
  role = aws_iam_role.events.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["ecs:RunTask", "iam:PassRole"]
      Resource = "*"
    }]
  })
}

resource "aws_cloudwatch_event_target" "run" {
  rule     = aws_cloudwatch_event_rule.schedule.name
  arn      = aws_ecs_cluster.this.arn
  role_arn = aws_iam_role.events.arn
  ecs_target {
    task_definition_arn = aws_ecs_task_definition.this.arn
    launch_type         = "FARGATE"
    network_configuration {
      subnets          = var.subnet_ids
      security_groups  = var.security_group_ids
      assign_public_ip = true
    }
  }
}
