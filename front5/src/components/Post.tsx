import React from "react";
import {Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";

export function ComponentPost(props: { post: Post }) {
    const post = props.post;

    return (
        <div className={styles.post}>
            <PostUser post={post}/>
            <div className={styles.text}>{post.text}</div>
        </div>
    )
}
