import React from "react";
import {Link} from "./Link";
import {Post} from "../store2/types";

export function ComponentPost(props: { post: Post }) {
    const post = props.post;

    return (
        <div>
            From:
            <Link href={post.fromHref}>
                {post.fromName}
            </Link>

            {post.text}
            <button>Update name</button>
            <Link href={post.detailsURL}>Детали</Link>
        </div>
    )
}
