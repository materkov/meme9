import React, {useEffect} from "react";
import {Post, QueryParams} from "../types";
import {Post as PostComponent, PostQuery} from "./Post";
import {api} from "../api";

export function PostPage(props: { id: string }) {
    const [post, setPost] = React.useState<string | undefined>(undefined);
    useEffect(() => {
        const q: QueryParams = {
            node: {
                id: props.id,
                inner: {
                    onPost: PostQuery,
                }
            }
        }
        api(q).then(data => {
            if (data.node?.type !== "Post") {
                return
            }

            // @ts-ignore
            window.store[data.node.user.id] = data.node.user;
            // @ts-ignore
            data.node.user = data.node.user.id;
            // @ts-ignore
            window.store[data.node.id] = data.node;

            setPost("1");
        })
    }, []);

    return <>{post && <PostComponent postId={props.id}/>}</>
}
