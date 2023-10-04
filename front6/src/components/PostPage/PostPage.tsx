import React, {useEffect} from "react";
import {usePostPage} from "../../store/postPage";
import * as styles from "./PostPage.module.css";
import {Post} from "../Post/Post";
import {useResources} from "../../store/resources";

export function PostPage() {
    let postId = window.document.location.pathname.substring(7);
    const postPage = usePostPage()
    const resources = useResources();

    useEffect(() => {
        postPage.fetch(postId);
    }, []);

    const post = resources.posts[postId];
    const error = postPage.errors[postId];

    return (
        <div>
            {post && <Post post={post}/>}

            {!post && <div>Loading...</div>}

            {error && <div className={styles.error}>error</div>}
        </div>
    )
}
