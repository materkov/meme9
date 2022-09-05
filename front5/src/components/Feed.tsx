import React from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {Feed as FeedRenderer, Post, User} from "../store/types";

export function Feed(props: { data: FeedRenderer }) {
    const posts = props.data[0];

    return <>
        <Composer/>
        <br/>
        <a href="https://oauth.vk.com/authorize?client_id=7260220&response_type=code&redirect_uri=http://localhost:3000/vk-callback">Авторизация</a>
        <br/>
        {posts.map(post => <ComponentPost post={post} key={post.id}/>)}
    </>
}
