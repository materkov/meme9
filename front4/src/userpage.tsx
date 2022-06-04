import React, {useEffect} from "react";
import {QueryParams, User} from "./types";
import {api} from "./api";
import {Post, PostQuery} from "./components/post";

export function UserPage(props: { id: string }) {
    const [user, setUser] = React.useState<User | undefined>();

    useEffect(() => {
        const query = UserPageQuery(props.id);
        api(query).then(data => data.node?.type == "User" && setUser(data.node))
    }, [])

    return <>{user && <User user={user}/>}</>
}

export const UserPageQuery = (userId: string): QueryParams => ({
    node: {
        include: true,
        id: userId,
        inner: {
            onUser: {
                name: {include: true},
                posts: {
                    include: true,
                    inner: PostQuery,
                }
            },
        }
    }
})

export function User(props: { user: User }) {
    return <>
        Name: {props.user.name}
        <hr/>
        Posts:
        <br/>
        {props.user.posts?.map(post => <Post post={post}/>)}
    </>
}
