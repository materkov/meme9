import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";

export class UserPage extends React.Component<schema.UserPageRenderer> {
    render() {
        return (
            <div>
                <h1>User {this.props.id}</h1>
                <br/>
                <Link href={"/posts/" + this.props.lastPostId} onClick={() => {
                }}>
                    Latest post {this.props.lastPostId}
                </Link>
                <br/>
                Name: {this.props.name}
                <br/><br/>
                You are user: {this.props.currentUserId}
            </div>
        );
    }
}
