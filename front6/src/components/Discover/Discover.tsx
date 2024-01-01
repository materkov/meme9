import React from "react";
import * as types from "../../api/api";
import {FeedType, PostsListReq} from "../../api/api";
import {useInfiniteQuery, useQueryClient} from '@tanstack/react-query'
import {PostsList} from "../Post/PostsList";
import {useGlobals} from "../../store/globals";
import {Composer} from "./Composer";
import {getAllFromPosts} from "../../utils/postsList";
import {usePrefetch} from "../../utils/prefetch";

export function Discover() {
    const queryClient = useQueryClient();
    const globalState = useGlobals();
    const [discoverState, setDiscoverState] = React.useState(types.FeedType.DISCOVER);

    usePrefetch('__postsList', (data: any) => {
        queryClient.setQueryData(['discover', types.FeedType.DISCOVER], {
            pages: [data],
            pageParams: [''],
        });
        getAllFromPosts(queryClient, data.items);

    })

    const {data, status, hasNextPage, fetchNextPage} = useInfiniteQuery({
        queryKey: ['discover', discoverState],
        queryFn: ({pageParam}) => {
            const req = new PostsListReq();
            req.pageToken = pageParam;
            req.count = 10;
            req.type = discoverState;

            return types.postsList(req).then(res => {
                getAllFromPosts(queryClient, res.items);
                return res;
            })
        },
        initialPageParam: "",
        getNextPageParam: (lastPage) => lastPage.pageToken,
    });

    const switchType = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        setDiscoverState(discoverState === FeedType.FEED ? FeedType.DISCOVER : FeedType.FEED);
    };

    return <div>
        <h1>Discover</h1>

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
