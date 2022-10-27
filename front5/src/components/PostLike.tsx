import styles from "./PostLike.module.css";
import {Heart} from "./icons/Heart";
import React from "react";
import {HeartRed} from "./icons/HeartRed";
import {fetcher, queryClient} from "../store/fetcher";
import {api, PostLikeData, Viewer} from "../store/types";
import {useMutation, useQuery} from "@tanstack/react-query";
import classNames from "classnames";
import {PostLikers} from "./PostLikers";

export const PostLike = (props: { id: string }) => {
    const queryKey = ["/posts/" + props.id + "/liked?count=0"];
    const {data} = useQuery<PostLikeData>(queryKey, fetcher);
    const {data: viewer} = useQuery<Viewer>(["/viewer"], fetcher);
    const [likersVisible, setLikersVisible] = React.useState(false);

    const unlike = useMutation(
        () => (api("/postUnlike", {id: props.id})),
        {
            onSuccess: () => {
                queryClient.invalidateQueries(queryKey);
                queryClient.invalidateQueries(["/posts/" + props.id + "/liked?count=10"]);
            },
            onMutate: () => {
                const prevData = queryClient.getQueryData<PostLikeData>(queryKey);
                if (!prevData) return;

                const newData = {...prevData, isViewerLiked: false, totalCount: (prevData.totalCount || 0) - 1};
                queryClient.setQueryData(queryKey, newData);
            }

        })
    const like = useMutation(
        () => api("/postLike", {id: props.id}),
        {
            onSuccess: () => {
                queryClient.invalidateQueries(queryKey);
                queryClient.invalidateQueries(["/posts/" + props.id + "/liked?count=10"]);
            },
            onMutate: () => {
                const prevData = queryClient.getQueryData<PostLikeData>(queryKey);
                if (!prevData) return;

                const newData = {...prevData, isViewerLiked: true, totalCount: (prevData.totalCount || 0) + 1};
                queryClient.setQueryData(queryKey, newData);
            }
        }
    )

    const onClick = () => {
        if (!data || !viewer?.viewerId) return;

        data?.isViewerLiked ? unlike.mutate() : like.mutate();
    }

    if (!data) return null;

    return <div className={styles.likeBtn} onClick={onClick}
                onMouseEnter={() => setLikersVisible(true)}
                onMouseLeave={() => setLikersVisible(false)}
    >
        {data.isViewerLiked ?
            <HeartRed className={styles.likeIcon}/> :
            <Heart className={styles.likeIcon}/>
        }

        {data.totalCount > 0 &&
            <div className={styles.likeText}>{data.totalCount}</div>
        }

        {likersVisible &&
            <div className={classNames({
                [styles.likersPopup]: true,
                [styles.likersPopup__visible]: true,
            })}>
                <PostLikers id={props.id}/>
            </div>
        }
    </div>
}
