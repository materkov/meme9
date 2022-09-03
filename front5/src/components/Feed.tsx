import React from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {Feed as FeedRenderer, Post, User} from "../store/types";

export function Feed(props: { data: FeedRenderer }) {
    const postIds = props.data.posts || [];

    let posts: Array<{ post: Post, user?: User }> = [];
    for (let postId of postIds) {
        const post = props.data.nodes?.posts?.find(item => item.id === postId);
        if (!post) {
            continue;
        }

        const user = props.data.nodes?.users?.find(item => item.id == post.fromId);

        posts.push({post, user});
    }

    return <>
        <Composer/>
        <br/>
        <a href="https://oauth.vk.com/authorize?client_id=7260220&response_type=code&redirect_uri=http://localhost:3000/vk-callback">Авторизация</a>
        <br/>
        {posts.map(({post, user}) => <ComponentPost post={post} from={user} key={post.id}/>)}
    </>
}
