import React from "react";
import * as styles from "./Profile.module.css";
import {postsList, PostsListReq, SubscribeAction, User, usersFollow, usersList, usersSetStatus} from "../../api/api";
import {useGlobals} from "../../store/globals";
import {PostsList} from "../Post/PostsList";
import {useInfiniteQuery, useQuery, useQueryClient} from "@tanstack/react-query";
import {getAllFromPosts} from "../../utils/postsList";
import {usePrefetch} from "../../utils/prefetch";

export function Profile() {
    const userId = document.location.pathname.substring(7);
    const queryClient = useQueryClient();
    const globals = useGlobals();

    usePrefetch('__userPage', (data: any) => {
        queryClient.setQueryData(['users', userId], data.user);
        queryClient.setQueryData(['userPosts', userId], {
            pages: [data.posts],
            pageParams: [''],
        });
        getAllFromPosts(queryClient, data.posts.items);
    });

    const {data: user} = useQuery({
        queryKey: ['users', userId],
        queryFn: () => (
            usersList({userIds: [userId]}).then(resp => resp[0])
        )
    })

    const {data: userPosts, hasNextPage, fetchNextPage} = useInfiniteQuery({
        queryKey: ['userPosts', userId],
        queryFn: ({pageParam}) => {
            const req = new PostsListReq();
            req.count = 10;
            req.byUserId = userId;
            req.pageToken = pageParam;

            return postsList(req).then(r => {
                getAllFromPosts(queryClient, r.items);
                return r;
            })
        },
        initialPageParam: '',
        getNextPageParam: (lastPage) => lastPage.pageToken
    })

    const [status, setStatus] = React.useState("");


    const updateStatus = () => {
        usersSetStatus({status: status})
            .then(() => {
                queryClient.setQueryData(['users', userId], (oldData: User) => {
                    const copy = structuredClone(oldData) as User;
                    copy.status = status;
                    queryClient.setQueryData(['users', userId], copy);
                })
            });
    };

    const follow = () => {
        if (!user) {
            return;
        }

        usersFollow({
            targetId: userId,
            action: user.isFollowing ? SubscribeAction.UNFOLLOW : SubscribeAction.FOLLOW,
        }).then(() => {
            queryClient.setQueryData(['users', userId], (oldData: User) => {
                const copy = structuredClone(oldData) as User;
                copy.isFollowing = true;
                queryClient.setQueryData(['users', userId], copy);
            })
        });
    }

    if (!user) {
        return <div>Loading....</div>
    }

    return <div>
        <h1 className={styles.userName}>{user.name}</h1>
        <div>{user.status}</div>

        {globals.viewerId === user.id && <>
        <textarea placeholder="Your text status..." className={styles.statusInput} value={status}
                  onChange={e => setStatus(e.target.value)}></textarea>
            <button onClick={updateStatus}>Update</button>
        </>
        }

        {globals.viewerId && globals.viewerId !== user.id && <>
            {user.isFollowing ?
                <button onClick={follow}>Unfollow user</button> :
                <button onClick={follow}>Follow user</button>
            }

        </>}

        <hr/>

        {userPosts?.pages.map((page, i) => (
            <PostsList postIds={page.items.map(post => post.id)} key={i}/>
        ))}

        {hasNextPage && <button onClick={() => fetchNextPage()}>Load more posts...</button>}
    </div>
}


