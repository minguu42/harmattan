resource "terraform_data" "hello_lambda_build" {
  triggers_replace = [
    [for f in fileset("../lambdas/hello", "*.go") : filemd5("../lambdas/hello/${f}")],
    filemd5("../../go.mod"),
    filemd5("../../go.sum")
  ]
  provisioner "local-exec" {
    command     = "GOOS=linux GOARCH=arm64 go build -o ./infra/lambdas/hello/bootstrap ./infra/lambdas/hello"
    working_dir = "${path.module}/../.."
  }
}

data "archive_file" "hello_lambda" {
  type        = "zip"
  source_file = "${path.module}/../lambdas/hello/bootstrap"
  output_path = "${path.module}/../lambdas/hello/hello.zip"
  depends_on  = [terraform_data.hello_lambda_build]
}

resource "aws_lambda_function" "hello" {
  filename         = data.archive_file.hello_lambda.output_path
  function_name    = "${local.product}-${var.env}-hello"
  role             = aws_iam_role.hello_lambda.arn
  handler          = "bootstrap"
  source_code_hash = data.archive_file.hello_lambda.output_base64sha256
  runtime          = "provided.al2023"
  architectures    = ["arm64"]
  timeout          = 60
  depends_on       = [aws_cloudwatch_log_group.hello_lambda]
}

resource "aws_iam_role" "hello_lambda" {
  name = "${local.product}-${var.env}-hello-lambda"
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy" "hello_lambda" {
  role = aws_iam_role.hello_lambda.id
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "${aws_cloudwatch_log_group.hello_lambda.arn}:*"
      }
    ]
  })
}

resource "aws_cloudwatch_log_group" "hello_lambda" {
  name              = "/aws/lambda/${local.product}-${var.env}-hello"
  retention_in_days = 3
}
