import React from "react";
import * as styles from "./PostPage.module.css";
import {Post} from "../Post/Post";
import {useQuery} from "@tanstack/react-query";
import * as types from "../../api/api";

export function PostPage() {
    let postId = window.document.location.pathname.substring(7);

    const {data, isLoading, error} = useQuery({
        queryKey: ['post', postId],
        queryFn: () => (
            types.postsListById({
                id: postId,
            }).then(res => {
                return res;
            })
        )
    })

    return (
        <div>
            {data && <Post postId={postId}/>}

            {isLoading && <div>Loading...</div>}

            {error && <div className={styles.error}>error</div>}
        </div>
    )
}
