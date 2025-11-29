# キャッシュクリア手順

## 問題
`npm run dev`を実行しても、以前のアプリのUIが表示される場合の対処法

## 解決方法

### 1. Next.jsのビルドキャッシュを削除（推奨：まずこれを試す）

```bash
cd frontend
npm run clean
npm run dev
```

### 2. 完全クリーンアップ（上記で解決しない場合）

```bash
cd frontend
npm run clean:all
npm install
npm run dev
```

### 3. 手動でキャッシュを削除する場合

以下のフォルダ/ファイルを削除：
- `.next/` - Next.jsのビルドキャッシュ
- `node_modules/.cache/` - 依存パッケージのキャッシュ
- `node_modules/` - 完全に再インストールしたい場合
- `package-lock.json` - ロックファイルを再生成したい場合

### 4. ブラウザのキャッシュもクリア

#### Chrome/Edgeの場合：
1. 開発者ツールを開く（F12）
2. ネットワークタブを開く
3. 「キャッシュの無効化」にチェックを入れる
4. ページをリロード（Ctrl+Shift+R または Cmd+Shift+R）

または：
1. 設定 → プライバシーとセキュリティ → 閲覧履歴データの削除
2. 「キャッシュされた画像とファイル」を選択
3. 削除

#### Firefoxの場合：
1. 開発者ツールを開く（F12）
2. ネットワークタブを開く
3. 「キャッシュを無効化」にチェックを入れる
4. ページをリロード（Ctrl+Shift+R または Cmd+Shift+R）

### 5. 別のポートで起動する

```bash
npm run dev -- -p 3004
```

### 6. シークレット/プライベートモードで確認

ブラウザのシークレット/プライベートモードで開いて、キャッシュの影響を確認

## よくある原因

1. **Next.jsのビルドキャッシュ** - `.next`フォルダに残っている
2. **ブラウザのキャッシュ** - 古いJS/CSSファイルがキャッシュされている
3. **Service Worker** - 以前のアプリがService Workerを登録している場合
4. **ポートの競合** - 別のプロセスが同じポートを使っている

## 確認方法

```bash
# ポート3000が使用されているか確認
lsof -i :3000
# または
netstat -ano | findstr :3000  # Windows
netstat -an | grep 3000       # Linux/Mac

# プロセスを終了する場合
kill -9 <PID>  # Linux/Mac
taskkill /PID <PID> /F  # Windows
```

