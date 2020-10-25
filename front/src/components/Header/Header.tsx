import React from "react";
import {Link} from "../Link/Link";
import * as schema from "../../schema/login";

export class Header extends React.Component<schema.HeaderRenderer> {
    render() {
        return (
            <div>
                <h1>meme</h1>
                Your user ID: {this.props.currentUserId}<br/>
                <Link href={"/"}>Index</Link> | <Link href={"/feed"}>Feed</Link> | <Link href={"/login"}>Login</Link>

                <hr/><hr/>
            </div>
        );
    }
}
