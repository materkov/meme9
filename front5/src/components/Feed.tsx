import React from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {Feed as FeedRenderer, Post, User} from "../store/types";

export function Feed(props: { data: FeedRenderer }) {
    const posts = props.data[0];

    return <>
        <Composer/>
        {posts.map(post => <ComponentPost post={post} key={post.id}/>)}
    </>
}
