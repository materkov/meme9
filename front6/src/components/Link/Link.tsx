import React from "react";
import {navigationGo} from "../../store/navigation";

export function Link(props: {
    className?: string,
    href: string,
    children: React.ReactNode,
}) {
    const onClick = (e: React.MouseEvent) => {
        e.preventDefault();
        navigationGo(props.href);
    };

    return (
        <a href={props.href} className={props.className} onClick={onClick}>
            {props.children}
        </a>
    )
}