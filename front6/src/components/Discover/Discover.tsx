import React from "react";
import * as types from "../../api/api";
import {FeedType} from "../../api/api";
import {useInfiniteQuery, useQueryClient} from '@tanstack/react-query'
import {PostsList} from "../Post/PostsList";
import {useGlobals} from "../../store/globals";
import {Composer} from "./Composer";
import {getAllFromPosts} from "../../utils/postsList";

export function Discover() {
    const queryClient = useQueryClient();
    const globalState = useGlobals();
    const [discoverState, setDiscoverState] = React.useState(types.FeedType.DISCOVER);

    const {data, status, hasNextPage, fetchNextPage} = useInfiniteQuery({
        queryKey: ['discover', discoverState],
        queryFn: ({pageParam}) => (
            types.postsList({
                pageToken: pageParam,
                count: 10,
                type: discoverState,
            }).then(res => {
                getAllFromPosts(queryClient, res.items);
                return res;
            })
        ),
        initialPageParam: "",
        getNextPageParam: (lastPage) => lastPage.pageToken,
    });

    const switchType = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        setDiscoverState(discoverState === FeedType.FEED ? FeedType.DISCOVER : FeedType.FEED);
    };

    return <div>
        <h1>Discover22</h1>

        {globalState.viewerId && <Composer/>}

        {globalState.viewerId && <>
            This is {discoverState == FeedType.DISCOVER ? 'discover' : 'feed'}. <a href="#" onClick={switchType}>
            Switch to {discoverState == FeedType.DISCOVER ? 'feed' : 'discover'}
        </a>
        </>}

        {status == "success" && <>
            {data?.pages.map((page, i) => (
                <PostsList key={i} postIds={page.items.map(post => post.id)}/>
            ))}
        </>}

        {hasNextPage && <button onClick={() => fetchNextPage()}>Load more</button>}
    </div>
}
