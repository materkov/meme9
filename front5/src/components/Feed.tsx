import React from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import {feedStore} from "../store/Feed";
import {useQuery} from "../store/fetcher";
import {Edges, Viewer} from "../store/types";

export function Feed() {
    const {data} = useQuery<Edges>("/feed");
    const {data: viewer, isLoading: isViewerLoading} = useQuery<Viewer>("/viewer");
    if (isViewerLoading) return <div>Loading...</div>;

    return <>
        {viewer?.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        {!data && <div>Loading...</div>}
        {data &&
            <PostsList posts={data.items} onShowMore={feedStore.fetch} showMore={Boolean(feedStore.nextCursor)}
                       showMoreDisabled={feedStore.isLoading}
            />
        }
    </>
}
