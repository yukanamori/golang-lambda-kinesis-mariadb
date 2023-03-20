# README

+ AWS IoT CoreからKinesis Data Stream経由でデータを受信し、Golang Lambdaで処理してAWS RDSに保存するサンプルコード
+ テストデータのフォーマットは以下

    ``` json
    {
        "field1": "string", # string
        "field2": 100,      # int
        "field3": true      # bool
    }
    ```

+ DB接続情報を書き換えること
