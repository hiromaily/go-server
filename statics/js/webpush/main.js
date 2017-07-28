/*
*
*  Push Notifications codelab
*  Copyright 2015 Google Inc. All rights reserved.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      https://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License
*
*/

/* eslint-env browser, es6 */

'use strict';

// TODO: is it better to get from server side??
const applicationServerPublicKey = 'BPl9ehqB6P27pIOGrzxp4Z5trdMmn1yvHeqbb4g-Q5SHyUhmK3CsCe9gU1QgxfFvbIQcTx-sc1ldYuhIAzGHZgQ';

const pushButton = document.querySelector('.js-push-btn');

let isSubscribed = false;
let swRegistration = null;

function urlB64ToUint8Array(base64String) {
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

// Check user is existing or not
function initialiseUI() {
    // Register user
    pushButton.addEventListener('click', function() {
        pushButton.disabled = true;
        if (isSubscribed) {
            unsubscribeUser();
        } else {
            subscribeUser();
        }
    });

    // Set the initial subscription value
    swRegistration.pushManager.getSubscription()
        .then(function(subscription) {
            // subscription includes json data
            isSubscribed = !(subscription === null);

            // commented out
            //updateSubscriptionOnServer(subscription);

            if (isSubscribed) {
                console.log('User IS subscribed.');
            } else {
                console.log('User is NOT subscribed.');
            }

            updateBtn();
        });
}

// Register user
function subscribeUser() {
    const applicationServerKey = urlB64ToUint8Array(applicationServerPublicKey);
    // 1.grant permission to user for displaying notification
    // 2.send request of push service to get information for generation of PushSubscription
    swRegistration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: applicationServerKey
    })
        .then(function(subscription) {
            // subscription includes json data
            console.log('User is subscribed:', subscription);

            updateSubscriptionOnServer(subscription);

            isSubscribed = true;

            updateBtn();
        })
        .catch(function(err) {
            console.log('Failed to subscribe the user: ', err);
            updateBtn();
        });
}

// Cancel user
function unsubscribeUser() {
    swRegistration.pushManager.getSubscription()
        .then(function(subscription) {
            if (subscription) {
                return subscription.unsubscribe();
            }
        })
        .catch(function(error) {
            console.log('Error unsubscribing', error);
        })
        .then(function() {
            updateSubscriptionOnServer(null);

            console.log('User is unsubscribed.');
            isSubscribed = false;

            updateBtn();
        });
}

function updateSubscriptionOnServer(subscription) {
    console.log("subscription:", subscription);

    // subscription includes json data
    // TODO: Send subscription to application server
    var data = {
        subscription: subscription,
        applicationKeys: {
            public: applicationServerPublicKey
        },
        data: ""
    };
    // {
    //     "endpoint": "https://fcm.googleapis.com/fcm/send/f_MHxDW6_8k:APA91bH998LbUYtT5Yjrr8NJak1liATortCgdN6EkshW7USjt8dUUpnf0ugVtLO7t35BbOStkHMMfKLad_-ndIdaVmg8dh7DQTk3kpJij4iA3aTRJwYey6v2PVfk3dsgW1gjRHYSSWsP",
    //     "expirationTime": null,
    //     "keys": {
    //         "p256dh": "BLt6u6ZEiqBK4-acClEFLjPZ4fvAT2Oo0nkpgKsiOxVYh95s8ODcJy7HJxOKbDq1RYypuO5ZbIzcMBQMOiaQFag=",
    //         "auth": "AZsanxKT5TnP1crxsE17_w=="
    //     }
    // }

    console.log(JSON.stringify(data));

    console.log('call fetch()');
    fetch("/webpush",
        {
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json'
            },
            method: "POST",
            body: JSON.stringify(data)
        })
        .then(function(res){
            console.log(res);
            displaySubscription(subscription);
        })
        .catch(function(res){ console.log(res) })

    //
    //displaySubscription(subscription)
}

// display subscription (below code is not important)
function displaySubscription(subscription) {
    console.log("displaySubscription()")
    const subscriptionJson = document.querySelector('.js-subscription-json');
    const subscriptionDetails =
        document.querySelector('.js-subscription-details');

    if (subscription) {
        subscriptionJson.textContent = JSON.stringify(subscription);
        subscriptionDetails.classList.remove('is-invisible');
    } else {
        subscriptionDetails.classList.add('is-invisible');
    }
}

// enable button and change text on button
function updateBtn() {
    // when denied
    if (Notification.permission === 'denied') {
        pushButton.textContent = 'Push Messaging Blocked.';
        pushButton.disabled = true;
        updateSubscriptionOnServer(null);
        return;
    }

    if (isSubscribed) {
        pushButton.textContent = 'Disable Push Messaging';
    } else {
        pushButton.textContent = 'Enable Push Messaging';
    }

    pushButton.disabled = false;
}

// Check whether Service Worker and push message are supported on current browser.
// And then, register sw.js for Service worker when these are suported.
if ('serviceWorker' in navigator && 'PushManager' in window) {
    console.log('Service Worker and Push is supported');



    navigator.serviceWorker.register('js/webpush/sw.js?ts=' + Date.now())
        .then(function(swReg) {
            console.log('Service Worker is registered', swReg);

            swRegistration = swReg;
            //
            initialiseUI();
        })
        .catch(function(error) {
            console.error('Service Worker Error', error);
        });
} else {
    console.warn('Push messaging is not supported');
    pushButton.textContent = 'Push Not Supported';
}
