import React from "react";
import {getByType, storeOnChanged} from "../store/store";

export function Link(props: { url: string, children: JSX.Element }) {
    const onClick = (e: React.MouseEvent) => {
        e.preventDefault();

        const routeObj = getByType("CurrentRoute");
        if (routeObj && routeObj.type === "CurrentRoute") {
            routeObj.url = props.url;
            storeOnChanged();
        }
    };

    return <a href={props.url} onClick={onClick}>{props.children}</a>
}
