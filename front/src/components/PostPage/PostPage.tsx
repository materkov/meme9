import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";

export class PostPage extends React.Component<schema.PostPageRenderer> {
    render() {
        return (
            <div>
                <h1>Post {this.props.id}</h1>
                {this.props.text}
                <br/>
                <img src="/static/cat.jpg" style={{width: "100px"}}/><br/>
                <Link href={"/users/" + this.props.userId} onClick={() => {}}>User {this.props.userId}</Link>
                <br/><br/>
                You are user: {this.props.currentUserId}
            </div>
        );
    }
}
