import React from "react";
import {Post} from "./Post/Post";
import {Composer} from "../Composer/Composer";
import * as schema from "../../api/api2";
import styles from "./Feed.module.css";

export function Feed(props: { data: schema.FeedRenderer }) {
    return <>
        {!props.data.placeholderText &&
        <Composer/>
        }

        {props.data.placeholderText &&
        <div className={styles.PlaceholderText}>{props.data.placeholderText}</div>
        }

        <PostsList posts={props.data.posts}/>
    </>
}

function PostsList(props: { posts: schema.Post[] }) {
    return <>
        {props.posts.map(post => (
            <Post key={post.id} data={post}/>
        ))}
    </>;
}
