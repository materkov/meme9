import React, {useEffect} from "react";
import {QueryParams, User} from "../types";
import {api, apiUpload} from "../api";
import {Post, PostQuery} from "./Post";
import styles from "./UserPage.module.css";

export function UserPage(props: { id: string }) {
    const [user, setUser] = React.useState<User | undefined>();
    const [viewer, setViewer] = React.useState<User | undefined>();

    useEffect(() => {
        const query = UserPageQuery(props.id);
        api(query).then(data => {
            if (data.node?.type == "User") {
                setUser(data.node)
            }

            setViewer(data.viewer);
        })
    }, [])

    return <>{user && <User user={user} viewer={viewer}/>}</>
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
    user: User;
    viewer?: User;
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
                        userId: props.user.id,
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
                        userId: props.user.id,
                    }
                }
            }
        };
        api(q).then(result => {
            location.reload();
        });
    }


    return <>
        Name: {props.user.name}

        <input type="file" className={styles.avatarInput} ref={avatarInputRef} onChange={onAvatarSelected}/>

        {props.viewer?.id === props.user.id &&
            <a className={styles.avatarUploadLabel} href={"#"} onClick={onChangeAvatar}>Сменить аватарку</a>
        }

        {props.user.id != props.viewer?.id &&
            <>
                {props.user.isFollowing && <div>Вы подписаны. <a href={"#"} onClick={onUnfollow}>Отписаться</a></div>}
                {!props.user.isFollowing && <div><a href={"#"} onClick={onFollow}>Подписаться</a></div>}
            </>
        }

        <hr/>
        {props.user.posts?.edges?.map(post => <Post key={post.id} post={post}/>)}
    </>
}
