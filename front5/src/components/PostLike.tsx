import styles from "./PostLike.module.css";
import {Heart} from "./icons/Heart";
import React from "react";
import {HeartRed} from "./icons/HeartRed";

export type Props = {
    count: number;
    isLiked: boolean;
    onToggle?: () => void;
}

export const PostLike = (props: Props) => {
    return <div className={styles.likeBtn} onClick={props.onToggle}>
        {props.isLiked ?
            <HeartRed className={styles.likeIcon}/> :
            <Heart className={styles.likeIcon}/>
        }

        {props.count > 0 &&
            <div className={styles.likeText}>{props.count}</div>
        }
    </div>
}
