import {Post as ApiPost} from "../../api/api";
import {Post} from "./Post";
import React from "react";
import * as styles from "./PostsLink.module.css";

export function PostsList(props: {
    postIds: string[]
}) {
    return <>
        {props.postIds.map((postId, idx) => <div key={postId}>
            <div className={styles.post}>
                <Post postId={postId}/>
            </div>

            {idx != props.postIds.length - 1 &&
                <div className={styles.postSeparator}/>
            }
        </div>)}
    </>

}