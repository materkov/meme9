import React, {useEffect} from "react";
import {vkCallback} from "../store/actions/auth";
import {setRoute} from "../store/actions/route";

export function VkCallback() {
    useEffect(() => {
        const urlParams = new URLSearchParams(location.search.substring(1));
        const code = urlParams.get("code") || "";

        vkCallback(code).then(() => {
            setRoute("/");
        })
    }, [])

    return <>Loading...</>;
}
