import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export interface UserPageProps {
    data: schema.UserPageRenderer;
}

export class UserPage extends React.Component<UserPageProps> {
    render() {
        return (
            <div>
                {this.props.data.headerRenderer && <Header data={this.props.data.headerRenderer}/>}

                <h1>User {this.props.data.id}</h1>
                <br/>
                <Link href={this.props.data.lastPostUrl} onClick={() => {
                }}>
                    Latest post {this.props.data.lastPostId}
                </Link>
                <br/>
                Name: {this.props.data.name}
                <br/><br/>
                You are user: {this.props.data.headerRenderer?.currentUserId}
            </div>
        );
    }
}
