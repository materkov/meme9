import React from "react";
import {Link} from "./Link";
import {Post, User} from "../store/types";
import {PostUser} from "./PostUser";

export function ComponentPost(props: { post: Post }) {
    const post = props.post;

    return (
        <div>
            {post.from && <PostUser user={post.from}/>}<br/>

            {post.text}
            <br/>
            <Link href={post.detailsURL}>Детали</Link>
            <br/>
            <br/>
            <br/>
        </div>
    )
}
