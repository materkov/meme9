import styles from "./PostLike.module.css";
import {Heart} from "./icons/Heart";
import React from "react";
import {HeartRed} from "./icons/HeartRed";
import {useQuery} from "../store/fetcher";
import {PostLikeData} from "../store/types";

export const PostLike = (props: {id: string}) => {
    const {data} = useQuery<PostLikeData>("/posts/" + props.id + "/isLiked");
    if (!data) return null;

    return <div className={styles.likeBtn}>
        {data.isLiked ?
            <HeartRed className={styles.likeIcon}/> :
            <Heart className={styles.likeIcon}/>
        }

        {data.likesCount > 0 &&
            <div className={styles.likeText}>{data.likesCount}</div>
        }
    </div>
}
