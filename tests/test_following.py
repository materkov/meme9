import http.client
import json
import urllib.parse
import pytest

conn = http.client.HTTPSConnection("meme.mmaks.me")
#conn = http.client.HTTPConnection("127.0.0.1", port=8000)

token = ""
user1_id = ""
user2_id = ""


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

    print('<--', method, params)
    print('-->', resp)
    print('')

    resp =  json.loads(resp)
    return resp.get('data', None), resp.get('error', None)


# user1: test@test.ru
# user2: test2@test.ru


def auth():
    global token, user1_id, user2_id
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


def test_follow():
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


def test_posting():
    post, err = api("posts.add", {
        "text": "Test post text"
    })
    assert err == None
    assert post["id"]
    assert post["text"] == "Test post text"
    assert post["userId"] == user1_id

    like_data, err = api("posts.like", {
        "postId": post["id"]
    })
    assert err == None
    assert like_data['totalCount'] == 1
    assert like_data['isViewerLiked'] == True

    like_data, err = api("posts.unlike", {
        "postId": post["id"]
    })
    assert err == None
    assert not like_data.get('totalCount', 0)
    assert not like_data.get('isViewerLiked', False)

    resp, err = api("posts.delete", {
        "id": post["id"]
    })
    assert err == None


auth()
test_follow()
test_posting()
