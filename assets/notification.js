
self.addEventListener('load', async () => {
  if ('serviceWorker' in navigator) {
    await navigator.serviceWorker
      .register('/sw.js', { scope: '/' })
      .then(function (reg) {
        window.sw = reg
        console.log('登録に成功しました。 Scope は ' + reg.scope);
      }).catch(function (error) {
        console.log('登録に失敗しました。' + error);
      });
  }
});
async function allowWebPush() {
  if ('Notification' in window) {
      let permission = Notification.permission;

      if (permission === 'denied') {
          alert('Push通知が拒否されているようです。ブラウザの設定からPush通知を有効化してください');
          return false;
      } else if (permission === 'granted') {
          alert('すでにWebPushを許可済みです');
      }
  }
  // 取得したPublicKey
  const appServerKey = 'BOUYK-PePdbLY7wGINbPsKGDY3LP0K7nB4XSkb8cNNi0PTCGzKCre88YzD0TWoVFVbmmPZjpFuIKoBMNy_ng7jU';
  const applicationServerKey = urlB64ToUint8Array(appServerKey);

  // push managerにサーバーキーを渡し、トークンを取得
  let subscription = undefined;
  try {
      subscription = await window.sw.pushManager.subscribe({
          userVisibleOnly: true,
          applicationServerKey
      });
  } catch (e) {
      alert('Push通知機能が拒否されたか、エラーが発生しましたので、Push通知は送信されません。');
      return false;
  }


  // 必要なトークンを変換して取得（これが重要！！！）
  const key = subscription.getKey('p256dh');
  const token = subscription.getKey('auth');
  const request = {
      endpoint: subscription.endpoint,
      userPublicKey: btoa(String.fromCharCode.apply(null, new Uint8Array(key))),
      userAuthToken: btoa(String.fromCharCode.apply(null, new Uint8Array(token)))
  };

  console.log(request);
  $.ajax({
    　url: '/event/subscribe', //アクセスするURLかディレクトリ
    　type: 'post',
    　cache: false, //cacheを使うかどうか
    　dataType:'json', //data type scriptなどデータタイプの指定
    　data: request, //アクセスするときに必要なデータを 記載
    　})
    　.done(function(response) { //通信が成功したときのコールバックの処理を書く
    　})
    　.fail(function(xhr) { //通信が失敗したときのコールバックの処理を書く
    　})
    　.always(function(xhr, msg) { //通信結果にかかわらず実行する処理を書く
    });
}



/**
* トークンを変換するときに使うロジック
* @param {*} base64String 
*/
function urlB64ToUint8Array (base64String) {
  const padding = '='.repeat((4 - base64String.length % 4) % 4);
  const base64 = (base64String + padding)
      .replace(/\-/g, '+')
      .replace(/_/g, '/');

  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
