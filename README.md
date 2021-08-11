# Main
This is simple docker image for sending message in given template (go template) format into Telegram.

# Inputs
Available input variables:
* `telegram_bot_token` - Telegram api bot token (required).
* `telegram_chat_id` - Telegram `chat_id` value (required).
* `env_var_prefix` - Environment variable names prefix. Default value: `TG_MSG_`. This environment variables may be used in message template.
* `template_string` - Message template string. If not defined, `template_base64_string` should be specified.
* `template_base64_string` - Base64 encoded message template string. If not defined, `template_string` should be specified.

# Use case
You can send notification messages, based on workflow status. For example:
```
name: notification test workflow

on:
  push:
    branches:
      - master

env:
  TELEGRAM_MESSAGE_BASE64_TEMPLATE: ${{ secrets.TELEGRAM_MESSAGE_BASE64_TEMPLATE }}
  TELEGRAM_BOT_TOKEN: ${{ secrets.TELEGRAM_BOT_TOKEN }}
  TELEGRAM_CHAT_ID: ${{ secrets.TELEGRAM_CHAT_ID }}

jobs:
  tests:
    runs-on: self-hosted
    steps:
      - uses: actions/checkout@v2

      - name: Get git commit info
        id: git_info
        run: |
          echo ::set-output name=version::$(git rev-parse --tags)
          echo ::set-output name=commit_id::${GITHUB_SHA:0:7}
          echo ::set-output name=date::$(git show -s --format=%cI HEAD)
          echo ::set-output name=message::$(git log --format=%B -n 1 HEAD)

      - name: Save variables to $GITHUB_ENV
        run: |
          echo TG_MSG_REPO_URL_PREFIX="https://github.com" >> $GITHUB_ENV
          echo TG_MSG_REPO_NAME="${{ github.repository }}" >> $GITHUB_ENV
          echo TG_MSG_WORKFLOW_NAME="${{ github.workflow }}" >> $GITHUB_ENV
          echo TG_MSG_WORKFLOW_RUN_ID="${{ github.run_id }}" >> $GITHUB_ENV
          echo TG_MSG_WORKFLOW_RUN_NUMBER="${{ github.run_number }}" >> $GITHUB_ENV
          echo TG_MSG_COMMIT_ID="${{ steps.git_info.outputs.commit_id }}" >> $GITHUB_ENV
          echo TG_MSG_BRANCH_NAME="${{ github.ref }}" >> $GITHUB_ENV
          echo TG_MSG_COMMIT_MESSAGE="${{ steps.git_info.outputs.message }}" >> $GITHUB_ENV


      - id: setup
        run: echo "some success step"

      - id: build
        run: |
          echo ::set-output name=status::failure
          exit 1

      - id: deploy
        run: |
          echo "second failed step"
          exit 1

      - if: success()
        uses: sputnik-systems/tg-message-sender@v0.0.2
        with:
          telegram_bot_token: ${{ env.TELEGRAM_BOT_TOKEN }}
          telegram_chat_id: ${{ env.TELEGRAM_CHAT_ID }}
          template_base64_string: ${{ env.TELEGRAM_MESSAGE_BASE64_TEMPLATE }}
        env:
          TG_MSG_WORKFLOW_OUTCOME: success
          TG_MSG_TAGS: test,success

      - if: failure()
        uses: sputnik-systems/tg-message-sender@v0.0.2
        with:
          telegram_bot_token: ${{ env.TELEGRAM_BOT_TOKEN }}
          telegram_chat_id: ${{ env.TELEGRAM_CHAT_ID }}
          template_base64_string: ${{ env.TELEGRAM_MESSAGE_BASE64_TEMPLATE }}
        env:
          TG_MSG_WORKFLOW_OUTCOME: failed
          TG_MSG_TAGS: test,failed
```
Before running it, you should write secrets into repository settings or org. level secrets. In this example, need define:
* `TELEGRAM_MESSAGE_BASE64_TEMPLATE` - base64 encoded message go template
* `TELEGRAM_BOT_TOKEN` - Telegram api bot token
* `TELEGRAM_CHAT_ID` - Telegram `chat_id` value

With this workflow example you can use template like this:
```
{{- $url_prefix := .TG_MSG_REPO_URL_PREFIX }}
{{- $repo_name := .TG_MSG_REPO_NAME }}
{{- $run_id := .TG_MSG_WORKFLOW_RUN_ID }}
{{- $tags := split .TG_MSG_TAGS "," }}
{{- $job_name := .TG_MSG_WORKFLOW_NAME }}
{{- $run_number := .TG_MSG_WORKFLOW_RUN_NUMBER }}
{{- $commit_id := .TG_MSG_COMMIT_ID }}
{{- $branch := .TG_MSG_BRANCH_NAME }}
{{- $status := .TG_MSG_WORKFLOW_OUTCOME }}
{{- $message := .TG_MSG_COMMIT_MESSAGE }}
{{ range $tags }}#{{ . }} {{ end }}
<b>job_name</b>: {{ $job_name }}
<b>build_id</b>: <a href="{{ $url_prefix }}/{{ $repo_name }}/actions/runs/{{ $run_id }}">{{ $run_number }}</a>
<b>commit_id</b>: <a href="{{ $url_prefix }}/{{ $repo_name }}/commit/{{ $commit_id }}">{{ $commit_id }}</a>
<b>branch</b>: {{ $branch }}
<b>status</b>: {{ $status }}

{{ $message }}
```
