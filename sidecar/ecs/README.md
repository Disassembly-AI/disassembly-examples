# AWS ECS (Fargate) sidecar

Two containers in one task: your `app` and a non-essential `disassembly` sidecar that scans it over
`localhost` and pulls `DISASSEMBLY_API_KEY` from Secrets Manager. Register with
`aws ecs register-task-definition --cli-input-json file://task-definition.json` and replace the
`<ACCOUNT_ID>` / `<REGION>` placeholders. For scheduled scans, drive this task from EventBridge
(see [`terraform/aws`](../../terraform/aws)).
