import React from "react";
import * as styles from "./Post.module.css";
import {Link} from "../Link/Link";
import {useGlobals} from "../../store/globals";
import * as types from "../../api/api";
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

export function Post(props: { postId: string }) {
    const {data: post, status} = useQuery<types.Post>({
        queryKey: ['post', props.postId],
    });

    if (status != 'success') {
        return null;
    }

    const globals = useGlobals();
    const queryClient = useQueryClient();

    const onDelete = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        types.postsDelete({postId: post.id}).then(() => {
            queryClient.invalidateQueries({queryKey: ['discover']});
            queryClient.invalidateQueries({queryKey: ['userPosts', post.user?.id]});
        });
    };

    const onLike = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        types.postsLike({
            postId: post.id,
            action: post.isLiked ? types.LikeAction.UNLIKE : types.LikeAction.LIKE,
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

    const onBookmark = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        if (!post) return;

        const cb = () => {
            queryClient.setQueryData(
                ['post', post.id],
                (oldData: types.Post) => {
                    const copy = structuredClone(oldData) as types.Post;
                    copy.isBookmarked = !copy.isBookmarked;

                    queryClient.setQueryData(['post', post.id], copy);
                }
            );

            queryClient.invalidateQueries({queryKey: ['bookmarks']});
        }

        if (post.isBookmarked) {
            types.bookmarksRemove({postId: post.id}).then(cb);
        } else {
            types.bookmarksAdd({postId: post.id}).then(cb);
        }
    };

    const date = new Date(post.date).toLocaleString();

    return <div className={styles.post}>
        <Link href={"/users/" + post.user?.id} className={styles.name}>{post.user?.name}</Link>

        <Link href={"/posts/" + post.id} className={styles.date}>{date}</Link>

        <div>
            {nl2br(post.text)}
        </div>

        {post.link && <LinkAttach link={post.link}/>}

        {post.poll && <Poll pollId={post.poll.id}/>}

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

        <div className={styles.bookmarksLine}>
            {globals.viewerId &&
                <>
                    <a onClick={onBookmark} href="#">
                        {post.isBookmarked ? 'Remove from bookmarks' : 'Add to bookmarks'}
                    </a>
                </>
            }
        </div>

        {globals.viewerId && post.user?.id == globals.viewerId &&
            <a onClick={onDelete} href="#" className={styles.deleteLink}>Delete post</a>
        }
    </div>
}
