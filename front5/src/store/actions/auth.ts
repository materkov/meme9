import * as types from "../../api/types";
import {store} from "../store";
import {api} from "../../api/api";
import {AuthVkCallback} from "../../api/types";

function setAuth(auth: types.Authorization) {
    store.dispatch({type: 'users/set', user: auth.user})
    store.dispatch({type: 'auth/setToken', token: auth.token})
    store.dispatch({type: 'viewer/set', userId: auth.user.id})

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
    store.dispatch({type: 'auth/setToken', token: ''})
    store.dispatch({type: 'viewer/set', userId: ''})

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
    api("auth.viewer", {} as AuthVkCallback).then((u: types.User) => {
        if (u) {
            store.dispatch({type: "users/set", user: u})
            store.dispatch({type: "viewer/set", userId: u.id})
        } else {
            store.dispatch({type: "viewer/set", userId: ''})
        }
    })
}
