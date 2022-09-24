import React from "react";
import styles from "./UserAvatar.module.css";

export type Props = {
    url?: string;
    width: number;
}

export function UserAvatar(props: Props) {
    const style = {
        width: props.width + 'px',
        height: props.width + 'px',
    };

    return <div>
        {props.url ?
            <img src={props.url} className={styles.avatar} alt="" style={style}/> :
            <div className={styles.avatar} style={style}/>
        }
    </div>;
}