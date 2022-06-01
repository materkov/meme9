import React, {useEffect} from "react";
import {QueryParams, User} from "./types";
import {api} from "./api";

export function UserPage(props: { id: string }) {
    const [user, setUser] = React.useState<User | undefined>();

    useEffect(() => {
        const q: QueryParams = {
            node: {
                include: true,
                id: props.id,
                inner: {
                    onUser: {
                        name: {include: true}
                    },
                }
            }
        }
        api(q).then(data => data.node?.type == "User" && setUser(data.node))
    }, [])

    return <>{user && <User user={user}/>}</>
}

export function User(props: { user: User }) {
    return <>
        Name: {props.user.name}
    </>
}
