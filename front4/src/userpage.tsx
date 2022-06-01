import React, {useEffect} from "react";
import {QueryParams, User} from "./types";
import {api} from "./api";

export function UserPage(props: { id: string }) {
    const [user, setUser] = React.useState<User | undefined>();

    useEffect(() => {
        const q: QueryParams = {
            user: {
                include: true,
                id: props.id,
                inner: {
                    name: {include: true},
                }
            }
        }
        api(q).then(data => setUser(data.user))
    }, [])

    return <>{user && <User user={user}/>}</>
}

export function User(props: { user: User }) {
    return <>
        Name: {props.user.name}
    </>
}
