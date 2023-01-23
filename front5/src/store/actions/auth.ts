import * as types from "../../api/types";
import {store} from "../store";
import {SetUser} from "../reducers/users";
import {SetToken} from "../reducers/auth";
import {SetViewer} from "../reducers/viewer";
import {api} from "../../api/api";

function setAuth(auth: types.Authorization) {
    store.dispatch({type: 'users/set', user: auth.user} as SetUser)
    store.dispatch({type: 'auth/setToken', token: auth.token} as SetToken)
    store.dispatch({type: 'viewer/set', userId: auth.user.id} as SetViewer)

    localStorage.setItem("authToken", auth.token);
}

export function vkCallback(code: string): Promise<void> {
    return new Promise((resolve, reject) => {
        api("auth.vkCallback", {
            code: code,
            redirectUri: location.origin + location.pathname,
        } as types.AuthVkCallback).then((r: types.Authorization) => {
            setAuth(r);
            resolve();
        })
    })
}

export function logout() {
    store.dispatch({type: 'auth/setToken', token: ''} as SetToken)
    store.dispatch({type: 'viewer/set', userId: ''} as SetViewer)

    localStorage.removeItem("authToken");
}

export function emailLogin(req: types.AuthEmailLogin): Promise<void> {
    return new Promise((resolve, reject) => {
        api("auth.emailLogin", req).then((r: types.Authorization) => {
            setAuth(r);
            resolve();
        }).catch(err => {
            reject();
        })
    })
}

export function emailRegister(req: types.AuthEmailRegister): Promise<void> {
    return new Promise((resolve, reject) => {
        api("auth.emailRegister", req).then((r: types.Authorization) => {
            setAuth(r);
            resolve();
        }).catch(err => {
            reject();
        })
    })
}

export function loadViewer() {
    api("auth.viewer", {}).then((u: types.User) => {
        if (u) {
            store.dispatch({type: "users/set", user: u} as SetUser)
            store.dispatch({type: "viewer/set", userId: u.id} as SetViewer)
        } else {
            store.dispatch({type: "viewer/set", userId: ''} as SetViewer)
        }
    })
}
