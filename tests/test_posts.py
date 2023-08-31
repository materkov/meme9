import http.client
import json
import urllib.parse
import pytest

#conn = http.client.HTTPSConnection("meme.mmaks.me")
conn = http.client.HTTPConnection("127.0.0.1", port=8000)

token = ""

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

    if resp.status == 200:
        return json.loads(resp_body), None
    else:
        return None, resp_body.decode("utf-8")


def test_posting():
    _, err = api("posts.add", {
        "text": "",
    })
    assert err == "empty text"

    _, err = api("posts.add", {
        "text": "Hello",
    })
    assert err == "not authorized"


test_posting()
