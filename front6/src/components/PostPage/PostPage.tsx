import React from "react";
import * as styles from "./PostPage.module.css";
import {Post} from "../Post/Post";
import {useQuery, useQueryClient} from "@tanstack/react-query";
import * as types from "../../api/api";
import {PostsListReq} from "../../api/api";
import {getAllFromPosts} from "../../utils/postsList";
import {usePrefetch} from "../../utils/prefetch";

export function PostPage() {
    let postId = window.document.location.pathname.substring(7);
    const queryClient = useQueryClient();

    usePrefetch('__postPagePost', (data: any) => {
        queryClient.setQueryData(['post', postId], data);
        getAllFromPosts(queryClient, [data]);
    });

    const {data, isLoading, error} = useQuery({
        queryKey: ['post', postId],
        queryFn: () => {
            const req = new PostsListReq();
            req.byId = postId;

            return types.postsList(req).then(r => {
                getAllFromPosts(queryClient, r.items);
                return r.items[0];
            })
        }
    })

    return (
        <div>
            {data && <Post postId={postId}/>}

            {isLoading && <div>Loading...</div>}

            {error && <div className={styles.error}>error</div>}
        </div>
    )
}
