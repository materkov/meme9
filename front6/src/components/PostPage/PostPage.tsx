import React, {useEffect} from "react";
import {usePostPage} from "../../store/postPage";
import * as styles from "./PostPage.module.css";
import {Post} from "../Post/Post";

export function PostPage() {
    let postId = window.document.location.pathname.substring(7);
    const postPage = usePostPage()

    useEffect(() => {
        postPage.fetch(postId);
    }, []);

    const post = postPage.posts[postId];
    const error = postPage.errors[postId];

    return (
        <div>
            {post && <Post post={post}/>}

            {!post && <div>Loading...</div>}

            {error && <div className={styles.error}>error</div>}
        </div>
    )
}
