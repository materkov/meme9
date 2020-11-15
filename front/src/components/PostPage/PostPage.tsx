import React from "react";
import * as schema from "../../schema/api";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export interface PostPageProps {
    data: schema.PostPageRenderer;
}

export class PostPage extends React.Component<PostPageProps> {
    render() {
        return (
            <div>
                <Header data={this.props.data.headerRenderer}/>

                <h1>Post {this.props.data.id}</h1>
                {this.props.data.text}
                <br/>
                <img src="/static/cat.jpg" style={{width: "100px"}}/><br/>
                <Link href={this.props.data.userUrl}>User {this.props.data.userId}</Link>
            </div>
        );
    }
}
