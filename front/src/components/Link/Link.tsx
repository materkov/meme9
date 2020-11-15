import React from "react";
import {resolveRoute} from "../../RouteResolver";
import {fetchJs} from "../../JsFetcher";
import {NavigateContext} from "../../context";

interface Props {
    href: string;
    onClick?: () => void;
}

export class Link extends React.Component<Props, any> {
    static contextType = NavigateContext;

    onClick = (e: React.MouseEvent) => {
        e.preventDefault();

        if (this.props.onClick) {
            this.props.onClick();
        } else {
            this.context(this.props.href);
        }
    };

    onMouseEnter = () => {
        // Preload route info and JS code
        resolveRoute(this.props.href).then(r => {
            fetchJs(r.js);
        });
    };

    render() {
        return <a href={this.props.href} onMouseEnter={this.onMouseEnter}
                  onClick={this.onClick}>{this.props.children}</a>;
    }
}
