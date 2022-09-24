import React, {useEffect} from "react";
import {api} from "../store/types";
import {authorize, navigate} from "../utils/localize";

export function VkCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(location.search.substring(1));

        api("/vkCallback", {
            code: urlParams.get("code") || "",
            redirectUri: location.origin + location.pathname,
        }).then(r => {
            navigate("/");
            authorize(r[0]);
        })
    }, [])

    return <>Loading...</>;
}
