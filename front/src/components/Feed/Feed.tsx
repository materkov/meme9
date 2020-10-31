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

                {posts.map(item => (<FeedItem key={item.id} data={item}/>))}
            </div>
        );
    }
}

interface FeedItemProps {
    data: schema.PostPageRenderer;
}

class FeedItem extends React.Component<FeedItemProps, any> {
    render() {
        const post = this.props.data;

        return (
            <div>
                Post <Link href={post.postUrl}>#{post.id}</Link><br/>
                From User <Link href={post.userUrl}>#{post.userId}</Link><br/>
                <br/>
                {post.text}
                <br/>
                <hr/>
            </div>
        );
    }
}