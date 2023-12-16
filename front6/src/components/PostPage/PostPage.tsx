import React from "react";
import * as styles from "./PostPage.module.css";
import {Post} from "../Post/Post";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import * as types from "../../api/api";
import {getAllFromPosts} from "../../utils/postsList";

export function PostPage() {
    let postId = window.document.location.pathname.substring(7);
    const queryClient = useQueryClient();

    const {data, isLoading, error} = useQuery({
        queryKey: ['post', postId],
        queryFn: () => (
            types.postsListById({
                id: postId,
            }).then(r => {
                getAllFromPosts(queryClient, [r]);
                return r;
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
