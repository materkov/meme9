import React from "react";
import {GlobalContext} from "../../Context";

export interface Props {
    href: string;
    className?: string;
}

export class Link extends React.PureComponent<Props> {
    static contextType = GlobalContext;

    onClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        this.context(this.props.href);
    };

    render() {
        return (
            <a href={this.props.href} onClick={this.onClick} className={this.props.className}>
                {this.props.children}
            </a>
        );
    }
}
