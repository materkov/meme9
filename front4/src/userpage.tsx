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
        id: userId,
        inner: {
            onUser: {
                name: {},
                posts: {
                    edges: {
                        inner: PostQuery,
                    },
                }
            },
        }
    }
})

export function User(props: { user: User }) {
    return <>
        Name: {props.user.name}
        <hr/>
        {props.user.posts?.edges?.map(post => <Post key={post.id} post={post}/>)}
    </>
}
