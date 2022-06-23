import React, {useEffect} from "react";
import {Post, QueryParams} from "../types";
import {Post as PostComponent, PostQuery} from "./Post";
import {api} from "../api";

export function PostPage(props: { id: string }) {
    const [post, setPost] = React.useState<Post | undefined>(undefined);
    useEffect(() => {
        const q: QueryParams = {
            node: {
                id: props.id,
                inner: {
                    onPost: PostQuery,
                }
            }
        }
        api(q).then(data => data.node?.type === "Post" && setPost(data.node))
    }, []);

    return <>{post && <PostComponent post={post}/>}</>
}
