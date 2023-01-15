import * as types from "../../store/types";
import {api2, AuthEmailLogin, Authorization, AuthVkCallback} from "../../store/types";
import {store} from "../store";
import {SetToken, SetUser, SetViewer} from "../reducers";

function setAuth(auth: Authorization) {
    store.dispatch({type: 'users/set', user: auth.user} as SetUser)
    store.dispatch({type: 'auth/setToken', token: auth.token} as SetToken)
    store.dispatch({type: 'viewer/set', userId: auth.user.id} as SetViewer)

    localStorage.setItem("authToken", auth.token);
}

export function vkCallback(code: string): Promise<void> {
    return new Promise((resolve, reject) => {
        api2("auth.vkCallback", {
            code: code,
            redirectUri: location.origin + location.pathname,
        } as AuthVkCallback).then((r: types.Authorization) => {
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

export function emailLogin(req: AuthEmailLogin): Promise<void> {
    return new Promise((resolve, reject) => {
        api2("auth.emailLogin", req).then((r: Authorization) => {
            setAuth(r);
            resolve();
        }).catch(err => {
            reject();
        })
    })
}