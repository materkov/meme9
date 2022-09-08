import React, {MouseEvent, ReactNode} from "react";
import {emitCustomEvent} from "react-custom-events";

export function Link(props: { href?: string, children: ReactNode, className?: string }) {
    const onClick = (e: MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        window.history.pushState(null, '', props.href);

        //window.document.dispatchEvent(new Event('urlChanged'));
        emitCustomEvent('urlChanged');
    }

    return <a className={props.className} href={props.href} onClick={onClick}>{props.children}</a>
    }