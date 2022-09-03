import React from "react";
import {Link} from "./Link";
import {Post, User} from "../store2/types";

export function ComponentPost(props: { post: Post, from?: User }) {
    const post = props.post;

    return (
        <div>
            From:
            {props.from &&
                <Link href={props.from.href}>
                    {props.from.name}
                </Link>
            }

            {post.text}
            <button>Update name</button>
            <Link href={post.detailsURL}>Детали</Link>
        </div>
    )
}
