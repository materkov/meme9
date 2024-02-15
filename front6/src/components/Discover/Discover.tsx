import React, {useEffect} from "react";
import * as types from "../../api/api";
import {useInfiniteQuery, useQueryClient} from '@tanstack/react-query'
import {PostsList} from "../Post/PostsList";
import {useGlobals} from "../../store/globals";
import {Composer} from "./Composer";
import {getAllFromPosts} from "../../utils/postsList";
import {usePrefetch} from "../../utils/prefetch";
import {getEvents} from "../../utils/realtime";
import {ApiPosts} from "../../api/client";

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

    useEffect(() => {
        if (!globalState.viewerId) {
            return;
        }

        getEvents(globalState.viewerId, (data: any) => {
            if (data.type === "NEW_POST") {
                queryClient.invalidateQueries({queryKey: ['discover']});
            }
        })
    }, []);

    const {data, status, hasNextPage, fetchNextPage} = useInfiniteQuery({
        queryKey: ['discover', discoverState],
        queryFn: ({pageParam}) => {
            const req: types.ListReq = {
                pageToken: pageParam,
                count: 10,
                type: discoverState,
                byId: "",
                byUserId: "",
            };

            return ApiPosts.List(req).then(res => {
                getAllFromPosts(queryClient, res.items);
                return res;
            })
        },
        initialPageParam: "",
        getNextPageParam: (lastPage) => lastPage.pageToken,
    });

    const switchType = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        setDiscoverState(discoverState === types.FeedType.FEED ? types.FeedType.DISCOVER : types.FeedType.FEED);
    };

    return <div>
        <h1>Discover</h1>

        {globalState.viewerId && <Composer/>}

        {globalState.viewerId && <>
            This is {discoverState == types.FeedType.DISCOVER ? 'discover' : 'feed'}. <a href="#" onClick={switchType}>
            Switch to {discoverState == types.FeedType.DISCOVER ? 'feed' : 'discover'}
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
