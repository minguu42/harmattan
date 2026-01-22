resource "aws_scheduler_schedule" "hello" {
  name       = "${local.product}-${var.env}-hello-schedule"
  group_name = "default"
  flexible_time_window {
    mode = "OFF"
  }
  schedule_expression = "rate(15 minutes)"
  target {
    arn      = aws_lambda_function.hello.arn
    role_arn = aws_iam_role.scheduler.arn
  }
}

resource "aws_iam_role" "scheduler" {
  name = "${local.product}-${var.env}-scheduler"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "scheduler.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "scheduler_invoke_lambda" {
  role = aws_iam_role.scheduler.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect   = "Allow"
        Action   = "lambda:InvokeFunction"
        Resource = aws_lambda_function.hello.arn
      }
    ]
  })
}
