import {Post as ApiPost} from "../../api/api";
import {Post} from "./Post";
import React from "react";
import * as styles from "./PostsLink.module.css";

export function PostsList(props: {
    posts: ApiPost[]
}) {
    return <>
        {props.posts.map((post, idx) => <div key={post.id}>
            <div className={styles.post}>
                <Post post={post}/>
            </div>

            {idx != props.posts.length - 1 &&
                <div className={styles.postSeparator}/>
            }
        </div>)}
    </>

}