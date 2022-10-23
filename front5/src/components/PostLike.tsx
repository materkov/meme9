import styles from "./PostLike.module.css";
import {Heart} from "./icons/Heart";
import React from "react";
import {HeartRed} from "./icons/HeartRed";
import {queryClient, useQuery} from "../store/fetcher";
import {api, PostLikeData, Viewer} from "../store/types";
import {useMutation} from "@tanstack/react-query";

export const PostLike = (props: { id: string }) => {
    const queryKey = ["/posts/" + props.id + "/isLiked"];
    const {data} = useQuery<PostLikeData>("/posts/" + props.id + "/isLiked");
    const {data: viewer} = useQuery<Viewer>("/viewer");

    const unlike = useMutation(
        () => (api("/postUnlike", {id: props.id})),
        {
            onSuccess: () => {
                queryClient.invalidateQueries(queryKey);
            },
            onMutate: () => {
                const prevData = queryClient.getQueryData<PostLikeData>(queryKey);
                if (!prevData) return;

                const newData = {...prevData, isLiked: false, likesCount: (prevData.likesCount || 0) - 1};
                queryClient.setQueryData(queryKey, newData);
            }

        })
    const like = useMutation(
        () => api("/postLike", {id: props.id}),
        {
            onSuccess: () => {
                queryClient.invalidateQueries(queryKey);
            },
            onMutate: () => {
                const prevData = queryClient.getQueryData<PostLikeData>(queryKey);
                if (!prevData) return;

                const newData = {...prevData, isLiked: true, likesCount: (prevData.likesCount || 0) + 1};
                queryClient.setQueryData(queryKey, newData);
            }
        }
    )

    const onClick = () => {
        if (!data || !viewer?.viewerId) return;

        data?.isLiked ? unlike.mutate() : like.mutate();
    }

    if (!data) return null;

    return <div className={styles.likeBtn} onClick={onClick}>
        {data.isLiked ?
            <HeartRed className={styles.likeIcon}/> :
            <Heart className={styles.likeIcon}/>
        }

        {data.likesCount > 0 &&
            <div className={styles.likeText}>{data.likesCount}</div>
        }
    </div>
}
