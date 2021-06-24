import React from "react";
import {GlobalStoreContext} from "../../Context";

export interface Props {
    href: string;
    className?: string;
}

export class Link extends React.PureComponent<Props> {
    static contextType = GlobalStoreContext;

    onClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        this.context.navigate(this.props.href);
    };

    render() {
        return (
            <a href={this.props.href} onClick={this.onClick} className={this.props.className}>
                {this.props.children}
            </a>
        );
    }
}
