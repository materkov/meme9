import React from "react";
import {PostParams} from "../types";
import styles from "./Post.module.css";
import {getByID} from "../store/store";

export type PostProps = {
    postId: string;
}

export const PostQuery: PostParams = {
    date: {},
    text: {},
    user: {
        inner: {
            name: {},
            avatar: {},
        }
    }
}

export function Post(props: PostProps) {
    const post = getByID(props.postId);
    if (post.type !== "Post") {
        return null;
    }

    const user = getByID(post.user || "");
    if (user.type !== "User") {
        return null;
    }

    const date = new Date(1000 * (post.date || 0));
    const dateFormatted = date.toISOString().slice(0, 19).replace("T", " ");

    return (
        <div className={styles.post}>
            <div className={styles.user}>
                <a href={"/users/" + user.id}>
                    <img alt="" className={styles.userAvatar}
                         src={user.avatar}
                    />
                    {user.name}
                </a>

                &nbsp;·&nbsp;
                <a href={"/posts/" + post.id} className={styles.postLink}>{dateFormatted}</a>
            </div>


            <div className={styles.postText}>{post.text}</div>
        </div>
    )
}
