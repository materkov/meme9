import React, {useEffect} from "react";
import {QueryParams, User} from "../types";
import {api, apiUpload} from "../api";
import {Post, PostQuery} from "./Post";
import styles from "./UserPage.module.css";
import {getByID} from "../store/store";

export function UserPage(props: { id: string }) {
    const [loaded, setLoaded] = React.useState<boolean>(false);

    useEffect(() => {
        const query = UserPageQuery(props.id);
        api(query).then(data => {
            // @ts-ignore
            window.store = {
                query: {
                    type: "Query",
                    // @ts-ignore
                    viewer: data.viewer.id,
                }
            };

            if (data.node?.type == "User") {
                const postIds = [];
                // @ts-ignore
                window.store[data.node.id] = data.node;
                for (let post of data.node.posts?.edges || []) {
                    // @ts-ignore
                    window.store[post.id] = post;
                    // @ts-ignore
                    window.store[post.id].user =window.store[post.id].user.id;
                    // @ts-ignore
                    postIds.push(post.id);
                }
                // @ts-ignore
                window.store[data.node.id].posts.edges = postIds;
            }

            //setViewer(data.viewer);
            // @ts-ignore
            //window.store[data.viewer.id] = data.viewer;
            setLoaded(true);
        })
    }, [])

    return <>{loaded && <User userId={props.id}/>}</>
}

export const UserPageQuery = (userId: string): QueryParams => ({
    node: {
        id: userId,
        inner: {
            onUser: {
                name: {},
                avatar: {},
                posts: {
                    edges: {
                        inner: PostQuery,
                    },
                },
                isFollowing: {},
            },
        }
    },
    viewer: {
        inner: {},
    }
})

export type UserProps = {
    userId: string;
};

export function User(props: UserProps) {
    const avatarInputRef = React.useRef<HTMLInputElement>(null);

    const onChangeAvatar = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        avatarInputRef.current?.click();
    };

    const onAvatarSelected = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files) {
            return;
        }

        const file = e.target.files[0];
        apiUpload(file).then((data) => {
            console.log(data);
        })
    };

    const onFollow = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        const q: QueryParams = {
            mutation: {
                inner: {
                    follow: {
                        userId: props.userId,
                    }
                }
            }
        };
        api(q).then(result => {
            location.reload();
        });
    }

    const onUnfollow = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();

        const q: QueryParams = {
            mutation: {
                inner: {
                    unfollow: {
                        userId: props.userId,
                    }
                }
            }
        };
        api(q).then(result => {
            location.reload();
        });
    }

    const user = getByID(props.userId);
    if (user.type !== "User") {
        return null;
    }

    const q = getByID("query");
    if (q.type !== "Query") {
        return null;
    }

    const viewerId = q.viewer;

    return <>
        Name: {user.name}

        <input type="file" className={styles.avatarInput} ref={avatarInputRef} onChange={onAvatarSelected}/>

        {viewerId === user.id &&
            <a className={styles.avatarUploadLabel} href={"#"} onClick={onChangeAvatar}>Сменить аватарку</a>
        }

        {user.id != viewerId &&
            <>
                {user.isFollowing && <div>Вы подписаны. <a href={"#"} onClick={onUnfollow}>Отписаться</a></div>}
                {!user.isFollowing && <div><a href={"#"} onClick={onFollow}>Подписаться</a></div>}
            </>
        }

        <hr/>
        {user.posts?.edges?.map(postId => <Post key={postId} postId={postId}/>)}
    </>
}
