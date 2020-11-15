import React from "react";
import * as schema from "../../schema/api";
import {Link} from "../Link/Link";
import {Header} from "../Header/Header";

export interface FeedProps {
    data: schema.GetFeedRenderer
}

export class Feed extends React.Component<FeedProps, any> {
    render() {
        const posts = this.props.data.posts;

        return (
            <div>
                <Header data={this.props.data?.headerRenderer}/>

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
                <Link href={post.postUrl}>Post #{post.id}</Link><br/>
                <Link href={post.userUrl}>User #{post.userId}</Link><br/>
                <br/>
                {post.text}
                <br/>
                <hr/>
            </div>
        );
    }
}