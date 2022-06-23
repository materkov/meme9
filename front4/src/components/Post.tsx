import React from "react";
import {Post, PostParams} from "../types";
import styles from "./Post.module.css";

export type PostProps = {
    post: Post;
}

export const PostQuery: PostParams = {
    date: {},
    text: {},
    user: {
        inner: {
            name: {},
        }
    }
}

export function Post(props: PostProps) {
    const date = new Date(1000 * (props.post.date || 0));
    const dateFormatted = date.toISOString().slice(0, 19).replace("T", " ");

    return (
        <div className={styles.post}>
            <div className={styles.user}>
                <a href={"/users/" + props.post.user?.id}>
                    <img alt="" className={styles.userAvatar}
                         src="https://avatars.githubusercontent.com/u/3899280?s=48&v=4"/>
                    {props.post.user?.name}
                </a>

                &nbsp;·&nbsp;
                <a href={"/posts/" + props.post.id} className={styles.postLink}>{dateFormatted}</a>
            </div>


            <div className={styles.postText}>{props.post.text}</div>
        </div>
    )
}
