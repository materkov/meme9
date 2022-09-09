import {useEffect} from "react";
import {apiHost} from "../store/types";
import {emitCustomEvent} from "react-custom-events";
import React from "react";

export function VkCallback() {
    useEffect(() => {
        const f = new FormData();
        f.set("redirectUri", location.origin + location.pathname);

        const urlParams = new URLSearchParams(location.search.substring(1));
        f.set("code", urlParams.get("code") || "");

        fetch(apiHost + "/vkCallback", {
            method: 'POST',
            body: f,
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                localStorage.setItem("authToken", r[0]);
                window.history.pushState(null, '', '/');
                emitCustomEvent('urlChanged');
                emitCustomEvent('onAuthorized');
            })
    }, [])

    return <>Loading...</>;
}
