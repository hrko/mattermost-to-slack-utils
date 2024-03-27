# mattermost-to-slack-utils
MattermostからSlackに移行するために個人的に作ったツール郡

## はじめに
これらのツールを使う前に、[Mattermostの公式ドキュメント](https://docs.mattermost.com/onboard/migrate-from-slack.html#download-file-attachments-and-email-addresses)を一読して、移行の手順を理解しておくことをお勧めします。

また、これらのツールは個人的な利用を想定しているため、バグがあるかもしれません。そのため、利用する際は自己責任でお願いします。インポートする前に、エクスポートデータを確認し、問題がないか確認してください。

動作確認は Ubuntu 22.04 で行っています。

## cmd/fix-channel-name
`mmetl` で変換したエクスポートデータのチャンネル名を修正するツール。
日本語のチャンネル名がチャンネルIDに変換されてしまう問題を修正する。

```sh
$ go run ./cmd/fix-channel-name channels.json mattermost_import.jsonl mattermost_import_fixed.jsonl
```

- `channels.json`: Slackのエクスポートデータ(.zip)の中にある `channels.json`
- `mattermost_import.jsonl`: `mmetl` で変換したMattermostのエクスポートデータ
- `mattermost_import_fixed.jsonl`: 修正後のMattermostのエクスポートデータ

## cmd/replace-username
ユーザー名を置換するツール。データ移行後にGitLabでのSSOによるログインに切り替えた際に、ユーザー名が変わる影響で、メンションが切れてしまうので、事前にユーザー名を置換しておくためのツール。

```sh
$ go run ./cmd/replace-username user-replace.json mattermost_import_fixed.jsonl mattermost_import_fixed_replaced.jsonl
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

## cmd/fix-attachments-filename
`mmetl` で変換したエクスポートデータの添付ファイル名を修正するツール。添付ファイルの名前が `ファイルID_ファイル名` になってしまう問題と、日本語等の非ASCII文字が消されてしまう問題を修正する。

```sh
$ go run ./cmd/fix-attachments-filename export-with-emails-and-attachments.zip data/bulk-export-attachments mattermost_import_fixed_replaced.jsonl mattermost_import_fixed_replaced_filename.jsonl
```

- `export-with-emails-and-attachments.zip`: Slackのエクスポートデータ(.zip)。ただし、[slack-advanced-exporter](https://github.com/grundleborg/slack-advanced-exporter/releases/)を使って添付ファイルを追加したもの。
- `data/bulk-export-attachments`: `mmetl` 変換したエクスポートデータの添付ファイルのディレクトリ
- `mattermost_import_fixed_replaced.jsonl`: `replace-username` で置換したMattermostのエクスポートデータ
- `mattermost_import_fixed_replaced_filename.jsonl`: 修正後のMattermostのエクスポートデータ
