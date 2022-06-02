import React, {useEffect} from "react";
import {Post, QueryParams} from "./types";
import {api} from "./api";
import {Post as PostComponent, PostQuery} from "./components/post";

export function PostPage(props: { id: string }) {
    const [post, setPost] = React.useState<Post | undefined>(undefined);
    useEffect(() => {
        const q: QueryParams = {
            node: {
                include: true,
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
