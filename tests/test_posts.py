import http.client
import json
import urllib.parse
import pytest

#conn = http.client.HTTPSConnection("meme.mmaks.me")
conn = http.client.HTTPConnection("127.0.0.1", port=8000)

token = ""
user_id = ""

def api(method, params=dict({})):
    params = json.dumps(params)
    headers = {
        "Content-type": "application/json",
        "Authorization": "Bearer " + token
    }

    url = "/api/" + method
    conn.request("POST", url, params, headers)
    resp = conn.getresponse()
    resp_body = resp.read()

    print('<--', method, params)
    print('-->', resp.status, resp_body.decode('utf-8'))
    print('')

    resp_body = json.loads(resp_body)

    if resp_body.get('error'):
        return None, resp_body['error']

    return resp_body, None


def test_posting():
    # posts.add
    post, err = api("posts.add", {
        "text": "test post",
    })
    assert err == None

    # posts.list
    posts_list, err = api("posts.list", {})
    assert err == None
    assert posts_list['items'][0]['id'] == post['id']

    # posts.list by user_id
    posts, err = api("posts.list", {
        "byUserId": post['userId'],
    })
    assert err == None
    assert posts['items'][0]['id'] == post['id']

    # posts.list by id
    postById, err = api("posts.list", {
        "byId": post['id'],
    })
    assert err == None
    assert postById['items'][0]['id'] == post['id']

    # posts.delete
    _, err = api("posts.delete", {
        "postId": post['id'],
    })
    assert err == None


def test_user():
    users, err = api("users.list", {"userIds": [user_id]})
    assert err == None
    assert users[0]['id'] == user_id


def test_auth():
    global token, user_id

    resp, err = api("auth.register", {
        "email": "test@email.com",
        "password": "123456",
    })
    if err == 'EmailAlreadyRegistered':
        pass
    else:
        assert err == None

    resp, err = api("auth.login", {
        "email": "test@email.com",
        "password": "123456",
    })
    assert err == None
    assert resp['token']
    assert resp['userId']

    token = resp['token']
    user_id = resp['userId']

test_auth()
test_posting()
test_user()
