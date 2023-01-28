import http.client
import json
import urllib.parse
import pytest

conn = http.client.HTTPSConnection("meme.mmaks.me")
#conn = http.client.HTTPConnection("127.0.0.1", port=8000)

token = ""


def api(method, params=dict({})):
    params = json.dumps(params)
    headers = {
        "Content-type": "application/json",
        "Authorization": "Bearer " + token
    }

    url = "/api2?method=" + method
    conn.request("POST", url, params, headers)
    resp = conn.getresponse()
    resp = resp.read()

    print('<--', url, params)
    print('-->', resp)
    print('')

    resp =  json.loads(resp)
    return resp.get('data', None), resp.get('error', None)


# user1: test@test.ru
# user2: test2@test.ru

def test_follow():
    global token
    login_resp, err = api("auth.emailLogin", {
        "email": "test@test.ru",
        "password": "1234"
    })
    assert err == None

    login_resp2, err = api("auth.emailLogin", {
        "email": "test2@test.ru",
        "password": "1234"
    })
    assert err == None

    user1_id = login_resp['user']['id']
    user2_id = login_resp2['user']['id']

    token = login_resp['token']
    assert token

    resp, err = api("auth.viewer")
    assert err == None
    assert resp['id'] == login_resp['user']['id']

    _, err = api("users.follow", {
        "userId": user2_id
    })
    assert err == None

    resp, err = api("users.following.list", {'userId': user1_id, 'count': 10})
    assert err == None
    assert resp['totalCount'] == 1

    resp, err = api("users.followers.list", {'userId': user2_id, 'count': 10})
    assert err == None
    assert resp['totalCount'] == 1

    _, err = api("users.unfollow", {"userId": user2_id})
    assert err == None

    resp, err = api("users.following.list", {'userId': user1_id, 'count': 10})
    assert err == None
    assert resp.get('totalCount', 0) == 0

    resp, err = api("users.followers.list", {'userId': user2_id, 'count': 10})
    assert err == None
    assert resp.get('totalCount', 0) == 0

test_follow()
