import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../schema/login";

export interface HeaderProps {
    data?: schema.HeaderRenderer;
}

export class Header extends React.Component<HeaderProps> {
    render() {
        const data = this.props.data;
        return (
            <div>
                <h1>meme</h1>
                {data?.currentUserName ?
                    <span>Вы вошли как: <b>{data.currentUserName}</b></span> :
                    <span>Вы не авторизованы.</span>
                }

                <br/>

                {data?.links.map((link) =>
                    <span key={link.url} style={{marginRight: '4px'}}>
                        <Link href={link.url}>{link.label}</Link>
                    </span>
                )}

                <hr/>
                <hr/>
            </div>
        );
    }
}
