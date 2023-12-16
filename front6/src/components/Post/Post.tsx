import React from "react";
import * as styles from "./Post.module.css";
import {Link} from "../Link/Link";
import {useGlobals} from "../../store/globals";
import * as types from "../../api/api";
import {LikeAction, Post as ApiPost, postsDelete, postsLike} from "../../api/api";
import {LinkAttach} from "./LinkAttach";
import {Poll} from "./Poll";
import {useQuery, useQueryClient} from "@tanstack/react-query";

const nl2br = (string: string) => {
    if (string) {
        return string.split("\n").map((item, key) => {
            return (
                <span key={key}>{item}<br/></span>
            );
        });
    }
};

export function Post(props: {
    postId: string,
}) {
    const {data, status} = useQuery<ApiPost>({
        queryKey: ['post', props.postId],
    });
    if (status != 'success') {
        return null;
    }

    const post = data;
    const date = new Date(post.date).toLocaleString();
    const globals = useGlobals();
    const queryClient = useQueryClient();

    const onDelete = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        postsDelete({postId: post.id}).then(() => {
            queryClient.invalidateQueries({queryKey: ['discover']});
            alert('Post deleted');
        });
    };

    const onLike = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        postsLike({
            postId: post.id,
            action: post.isLiked ? LikeAction.UNLIKE : LikeAction.LIKE,
        }).then(() => {
            queryClient.setQueryData(
                ['post', post.id],
                (oldData: types.Post) => {
                    const copy = structuredClone(oldData) as types.Post;

                    if (copy.isLiked) {
                        copy.isLiked = false;
                        copy.likesCount = (copy.likesCount || 0) - 1;
                    } else {
                        copy.isLiked = true;
                        copy.likesCount = (copy.likesCount || 0) + 1;
                    }

                    queryClient.setQueryData(['post', post.id], copy);
                }
            )
        });
    };

    return <div className={styles.post}>
        <Link href={"/users/" + post.user?.id} className={styles.name}>{post.user?.name}</Link>

        <Link href={"/posts/" + post.id} className={styles.date}>{date}</Link>

        <div>
            {nl2br(post.text)}
        </div>

        {post.link && <LinkAttach link={post.link}/>}

        {globals.viewerId && post.user?.id == globals.viewerId &&
            <a onClick={onDelete} href="#" className={styles.deleteLink}>Delete post</a>
        }

        <div className={styles.likesLine}>
            {globals.viewerId &&
                <>
                    <a onClick={onLike} href="#">
                        {post.isLiked ? 'Unlike' : 'Like'}
                    </a>
                    {post.likesCount > 0 && <> | </>}
                </>
            }

            {post.likesCount > 0 && <>{post.likesCount} like(s)</>}
        </div>

        {post.poll && <Poll pollId={post.poll.id}/>}
    </div>
}
