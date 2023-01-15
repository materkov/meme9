import React, {MouseEvent, ReactNode} from "react";
import {actions} from "../store2/actions";

export function Link(props: { href?: string, children: ReactNode, className?: string }) {
    const onClick = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        if (props.href) {
            actions.setRoute(props.href);
        }
    }

    return <a className={props.className} href={props.href} onClick={onClick}>{props.children}</a>
}