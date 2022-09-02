import {QueryParams} from "../types";
import {api} from "../api";
import {getByType, storeOnChanged} from "./store";

export function vkAuth() {
    const q: QueryParams = {
        mutation: {
            inner: {
                vkAuthCallback: {
                    url: location.href,
                }
            }
        }
    }
    api(q).then(result => {
        localStorage.setItem("authToken", result.mutation?.vkAuth?.token || '');

        const route = getByType("CurrentRoute")
        if (route && route.type == "CurrentRoute") {
            route.url = "/";
            history.pushState(null, '', '/');
            storeOnChanged();
        }
    })
}
