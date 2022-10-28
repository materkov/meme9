import React from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import {fetcher, useQuery} from "../store/fetcher";
import {Edges, Viewer} from "../store/types";
import {useInfiniteQuery} from "@tanstack/react-query";

export function Feed() {
    const {data, hasNextPage, fetchNextPage, isFetching} = useInfiniteQuery<Edges>(
        ["/feed"],
        ({pageParam = ""}) => fetcher({queryKey: ["/feed?cursor=" + pageParam]}),
        {
            getNextPageParam: (lastPage) => {
                return lastPage.nextCursor || undefined;
            }
        }
    );

    const {data: viewer, isLoading: isViewerLoading} = useQuery<Viewer>("/viewer");

    const posts: string[] = [];
    for (const page of data?.pages || []) {
        for (const post of page.items) {
            posts.push(post);
        }
    }

    if (isViewerLoading) return <div>Loading...</div>;

    return <>
        {viewer?.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        {!data && <div>Loading...</div>}
        {data &&
            <PostsList posts={posts} onShowMore={fetchNextPage} showMore={hasNextPage}
                       showMoreDisabled={isFetching}
            />
        }
    </>
}
