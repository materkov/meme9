import React from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import {fetcher, useQuery} from "../store/fetcher";
import {Edges, Viewer} from "../store/types";
import {useInfiniteQuery} from "@tanstack/react-query";
import classNames from "classnames";
import styles from "./Feed.module.css";

enum FeedType {
    FEED = "FEED",
    DISCOVER = "DISCOVER",
}

export function Feed() {
    const [feedType, setFeedType] = React.useState<FeedType>(FeedType.DISCOVER);

    const {data, hasNextPage, fetchNextPage, isFetching} = useInfiniteQuery<Edges>(
        ["/feed", feedType],
        ({pageParam = ""}) => fetcher({queryKey: ["/feed?feedType=" + feedType + "&cursor=" + pageParam]}),
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

        <div>Тип ленты:&nbsp;
            <a onClick={() => setFeedType(FeedType.FEED)} className={classNames({
                [styles.currentFeedType]: feedType == FeedType.FEED
            })}>Мои подписки</a> |&nbsp;
            <a onClick={() => setFeedType(FeedType.DISCOVER)} className={classNames({
                [styles.currentFeedType]: feedType == FeedType.DISCOVER
            })}>Дискавер</a>
        </div>

        {!data && <div>Loading...</div>}
        {data &&
            <PostsList posts={posts} onShowMore={fetchNextPage} showMore={hasNextPage}
                       showMoreDisabled={isFetching}
            />
        }
    </>
}
