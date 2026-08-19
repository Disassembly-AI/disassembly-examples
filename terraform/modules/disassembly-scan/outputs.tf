output "cluster_arn" {
  value = aws_ecs_cluster.this.arn
}
output "task_arn" {
  value = aws_ecs_task_definition.this.arn
}
output "schedule" {
  value = aws_cloudwatch_event_rule.schedule.schedule_expression
}
