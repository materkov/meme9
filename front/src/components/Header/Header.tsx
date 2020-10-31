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
                <Link href={"/"}>Index</Link> | <Link href={"/feed"}>Feed</Link> | <Link href={"/login"}>Login</Link> | <Link href={"/composer"}>Composer</Link>

                <hr/><hr/>
            </div>
        );
    }
}
