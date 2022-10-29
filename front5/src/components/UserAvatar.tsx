import React from "react";
import styles from "./UserAvatar.module.css";
import {useQuery} from "@tanstack/react-query";
import {Online, User} from "../store/types";
import {fetcher} from "../store/fetcher";

export type Props = {
    userId: string;
    width: number;
}

export function UserAvatar(props: Props) {
    const {data} = useQuery<User>([`/users/${props.userId}`], fetcher, {
        enabled: !!props.userId,
    })
    const {data: online} = useQuery<Online>([`/users/${props.userId}/online`], fetcher, {
        enabled: !!props.userId,
    })

    const style = {
        width: props.width + 'px',
        height: props.width + 'px',
    };

    let styleOnline = {
        right: '2px',
        bottom: '2px',
    }

    if (props.width == 100) {
        styleOnline = {
            right: '8px',
            bottom: '8px',
        }
    }

    return <div className={styles.container}>
        {data && data.avatar ?
            <img src={data.avatar} className={styles.avatar}
                 alt={"Avatar " + data.name} style={style}
            /> :
            <div className={styles.avatar} style={style}/>
        }

        {online?.isOnline && <div className={styles.online} style={styleOnline}/>}
    </div>;
}