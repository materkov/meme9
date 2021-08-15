import React from "react";
import {UrlContext} from "../context";

export const Link = (props: { href: string | undefined, children: any }) => {
    function onClick(e: React.MouseEvent<HTMLAnchorElement>, navigate: any) {
        e.preventDefault();
        navigate(props.href);
    }

    return (
        <UrlContext.Consumer>
            {({url, navigate}) => (
                <a href={props.href} onClick={(e) => onClick(e, navigate)}>
                    {props.children}
                </a>
            )}
        </UrlContext.Consumer>
    );
}
