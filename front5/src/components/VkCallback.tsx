import React, {useEffect} from "react";
import {api} from "../store/types";
import {emitCustomEvent} from "react-custom-events";

export function VkCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(location.search.substring(1));

        api("/vkCallback", {
            code: urlParams.get("code") || "",
            redirectUri: location.origin + location.pathname,
        }).then(r => {
            localStorage.setItem("authToken", r[0]);
            window.history.pushState(null, '', '/');
            emitCustomEvent('urlChanged');
            emitCustomEvent('onAuthorized');
        })
    }, [])

    return <>Loading...</>;
}
