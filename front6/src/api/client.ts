import {AuthClientImpl, BookmarksClientImpl, PollsClientImpl, PostsClientImpl, UsersClientImpl} from "./api";
import {cookieAuthToken, getCookie} from "../utils/cookie";

class TwirpRpc {
    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array> {
        return new Promise<Uint8Array>((resolve, reject) => {
            let headers: Record<string, string> = {};

            let token = getCookie(cookieAuthToken);
            if (token) {
                headers['authorization'] = 'Bearer ' + token;
            }

            // TODO think about this func
            fetch('/api/' + service + "/" + method, {
                credentials: 'omit',
                method: 'POST',
                body: JSON.stringify(data),
                headers: headers,
            })
                .then(r => {
                    if (!r.ok) {
                        reject('http error');
                    } else if (r.status !== 200) {
                        reject("incorrect http status " + r.status)
                    }

                    return r.text()
                })
                .then(r => {
                    try {
                        const resp = JSON.parse(r);
                        if (resp.error) {
                            reject(resp.error)
                        } else {
                            resolve(resp);
                        }
                    } catch (e) {
                        reject('cannot parse json');
                    }
                })
                .catch(reject);
        })
    }
}

const rpc = new TwirpRpc()

export const ApiBookmarks = new BookmarksClientImpl(rpc);
export const ApiPolls = new PollsClientImpl(rpc);
export const ApiPosts = new PostsClientImpl(rpc)
export const ApiAuth = new AuthClientImpl(rpc);
export const ApiUsers = new UsersClientImpl(rpc);
