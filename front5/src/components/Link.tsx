import React, {MouseEvent, ReactNode} from "react";
import {navigate} from "../utils/localize";

export function Link(props: { href?: string, children: ReactNode, className?: string }) {
    const onClick = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        if (props.href) {
            navigate(props.href);
        }
    }

    return <a className={props.className} href={props.href} onClick={onClick}>{props.children}</a>
}