import React, {MouseEvent} from "react";
import {StateContext} from "./StoreConnector";
//import {state} from "../state";

export const Link = (props: { href: string, children?: any }) => {
    const state = React.useContext(StateContext);
    const onClick = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        state.do({navigate: {url: props.href}})
    };

    return <a href={props.href} onClick={onClick}>{props.children}</a>
}
