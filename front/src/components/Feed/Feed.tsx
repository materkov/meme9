import React from "react";
import * as schema from "../../schema/login";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export class Feed extends React.Component<schema.GetFeedRenderer, any> {
    render() {
        const posts = this.props.posts;

        return (
            <div>
                {this.props.headerRenderer && <Header data={this.props.headerRenderer}/>}

                {posts.map(item => (<FeedItem key={item.id} post={item}/>))}
            </div>
        );
    }
}

class FeedItem extends React.Component<{ post: schema.PostPageRenderer }, any> {
    render() {
        const post = this.props.post;

        return (
            <div>
                Post <Link href={"/posts/" + post.id}>#{post.id}</Link><br/>
                From User <Link href={"/users/" + post.userId}>#{post.userId}</Link><br/>
                <br/>
                {post.text}
                <br/>
                <hr/>
            </div>
        );
    }
}