import React from "react";
import styles from "./UserAvatar.module.css";
import {useQuery} from "@tanstack/react-query";
import {User} from "../store/types";
import {fetcher} from "../store/fetcher";

export type Props = {
    userId: string;
    width: number;
}

export function UserAvatar(props: Props) {
    const {data} = useQuery<User>([`/users/${props.userId}`], fetcher, {
        enabled: !!props.userId,
    })
    const style = {
        width: props.width + 'px',
        height: props.width + 'px',
    };

    return <div>
        {data && data.avatar ?
            <img src={data.avatar} className={styles.avatar}
                 alt={"Avatar " + data.name} style={style}
            /> :
            <div className={styles.avatar} style={style}/>
        }
    </div>;
}