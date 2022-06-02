import React from "react";
import {Post, PostParams} from "../types";

export type PostProps = {
    post: Post;
}

export const PostQuery: PostParams = {
    date: {include: true},
    text: {include: true},
    user: {
        include: true,
        inner: {
            name: {include: true},
        }
    }
}

export function Post(props: PostProps) {
    return <div>
        <div><b>Text: </b> {props.post.text}</div>
        <div><b>User: </b> <a href={"/users/" + props.post.user?.id}>{props.post.user?.name}</a></div>
        <div><a href={"/posts/" + props.post.id}>Link</a></div>
        <hr/>
    </div>
}
