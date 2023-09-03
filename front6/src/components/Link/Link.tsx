import React from "react";
import {useNavigation} from "../../store/navigation";

export function Link(props: {
    className?: string,
    href: string,
    children: React.ReactNode,
}) {
    const navState = useNavigation();

    const onClick = (e: React.MouseEvent) => {
        e.preventDefault();
        navState.setURL(props.href);
    };

    return (
        <a href={props.href} className={props.className} onClick={onClick}>
            {props.children}
        </a>
    )
}