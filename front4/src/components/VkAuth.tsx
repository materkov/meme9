import React, {useEffect} from "react";
import {QueryParams} from "../types";
import {api} from "../api";
import {getByType, storeOnChanged} from "../store/store";

export function VKAuth() {
    useEffect(() => {
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
            //history.pushState(null, '', '/');
            //location.href = '/';

            const route = getByType("CurrentRoute")
            if (route && route.type == "CurrentRoute") {
                route.url = "/";
                history.pushState(null, '', '/');
                storeOnChanged();
            }
        })

    }, []);

    return <div>Авторизация...</div>;
}
