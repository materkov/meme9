import React from "react";
import * as schema from "../../schema/login";

export class PostPage extends React.Component<schema.PostPageRenderer> {
    render() {
        return (
            <div>
                <h1>Post {this.props.id}</h1>
                {this.props.text}
                <br/>
                <img src="/static/cat.jpg" style={{width: "100px"}}/><br/>
                <a href={"/users/" + this.props.userId}>User {this.props.userId}</a>
                <br/><br/>
                You are user: {this.props.currentUserId}
            </div>
        );
    }
}
