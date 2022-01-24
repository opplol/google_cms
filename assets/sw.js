self.addEventListener('install', (event) => {
  console.log('Installed');
});

self.addEventListener('activate', (event) => {
  console.log('Activated');
});

self.addEventListener('fetch', (event) => {
  console.log('Fetch request');
});

// プッシュ通知を受け取ったときのイベント
self.addEventListener('push', (event) => {
  console.log('Push event');
  const title = 'Push通知テスト';
  const options = {
      body: event.data.text(), // サーバーからのメッセージ
      tag: title, // タイトル
      icon: '', // アイコン
      badge: '' // アイコン
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

// プッシュ通知をクリックしたときのイベント
self.addEventListener('notificationclick', function (event) {
  event.notification.close();

  event.waitUntil(
      // プッシュ通知をクリックしたときにブラウザを起動して表示するURL
      clients.openWindow('http://localhost:8080/event/index')
  );
});
