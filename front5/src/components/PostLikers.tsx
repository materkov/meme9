import React from "react";
import styles from "./PostLikers.module.css";
import {Link} from "./Link";
import {useQueries, useQuery} from "@tanstack/react-query";
import {Edges, User} from "../store/types";
import {fetcher} from "../store/fetcher";
import {UserAvatar} from "./UserAvatar";

export function PostLikers(props: { id: string }) {
    const {data, isLoading, isStale} = useQuery<Edges>(["/posts/" + props.id + "/liked?count=10"], fetcher, {
    })

    const userQueries = useQueries<User[]>({
        queries: (data?.items || []).map(userId => {
            return {
                queryKey: ["/users/" + userId],
                queryFn: fetcher,
            }
        })
    })

    const users: User[] = [];
    for (let q of userQueries) {
        users.push(q.data as User);
    }

    return <div className={styles.list}>
        {(isLoading || isStale) && <>Загрузка...</>}
        {!isLoading && !isStale && !data?.totalCount && <>Никто не полайкал.</>}

        {!isLoading && !isStale && data && data.items?.map((userId, idx) => (
            <Link className={styles.item} href={"/users/" + userId} key={userId}>
                <UserAvatar width={40} userId={userId}/>
                <div className={styles.name}>{users[idx].name}</div>
            </Link>
        ))}
    </div>
}
