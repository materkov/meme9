import React, {useEffect} from "react";
import {vkCallback} from "../store2/actions/auth";
import {actions} from "../store2/actions";

export function VkCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(location.search.substring(1));
        const code = urlParams.get("code") || "";

        vkCallback(code).then(() => {
            actions.setRoute("/");
        })
    }, [])

    return <>Loading...</>;
}
