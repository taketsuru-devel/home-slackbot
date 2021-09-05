## 自宅のラズパイとか鍵をslackから操作したい

ひとまずゴリゴリ書いてたのをテストできるように軌道修正中

### memo

slackのevent subscriptionに登録したendpointはchallege処理が必要で関数化
bodyはjson形式
incoming webhookに登録したやつはchallenge処理不要
bodyはhttp form形式
