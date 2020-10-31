import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../schema/login";

export interface HeaderProps {
    data: schema.HeaderRenderer;
}

export class Header extends React.Component<HeaderProps> {
    render() {
        const data = this.props.data;
        return (
            <div>
                <h1>meme</h1>
                Вы вошли как: <b>{data.currentUserName}</b><br/>

                {data.links.map((link) =>
                    <Link key={link.url} href={link.url}>{link.label}</Link>
                )}

                <hr/>
                <hr/>
            </div>
        );
    }
}
