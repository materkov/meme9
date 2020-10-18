import React from "react";
import * as schema from "../../schema/login";

export class UserPage extends React.Component<schema.UserPageRenderer> {
    render() {
        return (
            <div>
                <h1>User {this.props.id}</h1>
                <br/>
                <a href={"/posts/" + this.props.lastPostId}>Latest post {this.props.lastPostId}</a>
                <br/>
                Name: {this.props.name}
                <br/><br/>
                You are user: {this.props.currentUserId}
            </div>
        );
    }
}
