# mattermost-to-slack-utils
MattermostからSlackに移行するために個人的に作ったツール

## cmd/fix-channel-name
`mmetl` で変換したエクスポートデータのチャンネル名を修正するツール。
日本語のチャンネル名がチャンネルIDに変換されてしまう問題を修正する。

```sh
$ go run cmd/fix-channel-name/main.go channels.json mattermost_import.jsonl mattermost_import_fixed.jsonl
```

- `channels.json`: Slackのエクスポートデータ(.zip)の中にある `channels.json`
- `mattermost_import.jsonl`: `mmetl` で変換したMattermostのエクスポートデータ
- `mattermost_import_fixed.jsonl`: 修正後のMattermostのエクスポートデータ

## cmd/replace-username
ユーザー名を置換するツール。データ移行後にGitLabでのSSOによるログインに切り替えた際に、ユーザー名が変わる影響で、メンションが切れてしまうので、事前にユーザー名を置換しておくためのツール。

```sh
$ go run cmd/replace-username/main.go user-replace.json mattermost_import_fixed.jsonl mattermost_import_fixed_replaced.jsonl
```

- `user-replace.json`: 置換前と置換後のユーザー名の対応表(形式は以下を参照)
- `mattermost_import_fixed.jsonl`: `fix-channel-name` で修正したMattermostのエクスポートデータ
- `mattermost_import_fixed_replaced.jsonl`: 置換後のMattermostのエクスポートデータ

`user-replace.json` の形式は以下の通り。

```json
[
{"name_old":"<元のユーザー名1>","name_new":"<新しいユーザー名1>"},
{"name_old":"<元のユーザー名2>","name_new":"<新しいユーザー名2>"},
...
{"name_old":"<元のユーザー名3>","name_new":"<新しいユーザー名3>"}
]
```
