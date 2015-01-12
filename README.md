# archer
## 構成
### ページ
#### channelを選択して時系列にログを見れる
* elasticsearchからchannelを指定してログを出す

#### ログをchannel, 文言で検索し、時間で絞り込めるページ

### API
#### ログを受けるAPI
Indexはslackに統一
typeはchannel毎に作る

#### mapping

```
curl -XPOST localhost:9200/slack -d @mapping.json
```

