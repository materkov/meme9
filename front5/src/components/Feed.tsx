import React from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {Feed as FeedRenderer, Post} from "../store2/types";

export function Feed(props: { data: FeedRenderer }) {

    props.data.nodes?.posts?.reduce(function(map, post) {});
    const postIds = props.data.posts;

    let posts = [];
    for (let postId of postIds) {
        const post = props.data.nodes?.posts?.find(item => item.id === postId);
        if (post) {
            posts.push(post);
        }
    }

    return <>
        <Composer/>
        {posts.map(post => <ComponentPost post={post} key={post.id}/>)}
    </>
}
