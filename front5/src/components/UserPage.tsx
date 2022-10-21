import React, {useEffect} from "react";
import {api, User, UserPostsConnection} from "../store/types";
import {ComponentPost} from "./Post";
import styles from "./UserPage.module.css";
import produce from "immer";
import {localizeCounter} from "../utils/localize";
import {UserAvatar} from "./UserAvatar";

export function UserPage() {
    const [user, setUser] = React.useState<User>();
    const [postsCursor, setPostsCursor] = React.useState("");
    const [viewerId, setViewerId] = React.useState("");
    const [loaded, setLoaded] = React.useState(false);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

    const [isFollowing, setIsFollowing] = React.useState(false);

    const [avatarUploading, setAvatarUploading] = React.useState(false);

    useEffect(() => {
        api("/userPage", {
            id: location.pathname.substring(7)
        }).then((r: [User, string]) => {
            const [user, viewerId] = r;
            setUser(user);
            setViewerId(viewerId);
            setUserName(user.name || "");

            // TODO strange
            for (let post of user.posts?.items || []) {
                post.user = user;
            }

            setLoaded(true);
            setPostsCursor(user.posts?.nextCursor || "");
            setIsFollowing(Boolean(user.isFollowing))
        })
    }, []);

    const editName = () => {
        api("/userEdit", {
            id: viewerId,
            name: userName,
        }).then(() => setUserNameUpdated(true));
    }

    if (!loaded || !user) {
        return <>Загрузка ...</>;
    }

    const onShowMore = () => {
        api("/userPage/posts", {
            id: location.pathname.substring(7),
            cursor: postsCursor,
        }).then((result: [UserPostsConnection]) => {
            let r = result[0];
            setPostsCursor(r.nextCursor || "");

            for (let post of r.items || []) {
                post.user = user;
            }

            setUser(produce(user, (user: User) => {
                user.posts = user.posts || {};
                user.posts.items = user.posts?.items || [];
                user.posts.items = [...user.posts.items, ...r.items || []];
            }));
        })
    }

    const onPostDelete = (postId: string) => {
        setUser(produce(user, user => {
            user.posts = user.posts || {};
            user.posts.items = user.posts?.items || [];
            user.posts.items = user.posts?.items?.filter(post => post.id !== postId);
        }))
    }

    const onFollow = () => {
        api("/userFollow", {
            id: user.id,
        }).then(() => {
            setIsFollowing(true);
        })
    }

    const onUnfollow = () => {
        api("/userUnfollow", {
            id: user.id,
        }).then(() => {
            setIsFollowing(false);
        })
    }

    const onAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files) {
            return;
        }

        setAvatarUploading(true);
        const file = e.target.files[0];

        api("/uploadAvatar", {
            file: file,
        }).then((resp) => {
            setUser({...user, avatar: resp.avatar});
            setAvatarUploading(false);
        })

    }

    return (
        <div>
            <div className={styles.topBlock}>
                <UserAvatar width={100} url={user.avatar}/>
                <div className={styles.infoBlock}>
                    <div className={styles.userName}>{user.name}</div>
                    <div className={styles.userBio}>
                        {user.bio}
                    </div>
                    <div className={styles.userCounters}>
                        <div className={styles.userCounter}>
                            <b>{user.posts?.count}</b> {localizeCounter(user.posts?.count || 0, "публикация", "публикации", "публикаций")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{user.followedByCount || 0}</b> {localizeCounter(user.followedByCount || 0, "подписчик", "подписчика", "подписчиков")}
                        </div>
                        <div className={styles.userCounter}>
                            <b>{user.followingCount || 0}</b> {localizeCounter(user.followingCount || 0, "подписка", "подписки", "подписок")}
                        </div>
                        <div className={styles.buttonsBlock}>
                            {viewerId && viewerId != user.id &&
                                <>
                                    {isFollowing ?
                                        <button onClick={onUnfollow}>Отписаться</button> :
                                        <button onClick={onFollow}>Подписаться</button>
                                    }
                                </>
                            }
                        </div>
                    </div>
                </div>
            </div>

            {user.id === viewerId && <>
                Имя: <input type="text" value={userName} onChange={e => setUserName(e.target.value)}/>
                <button onClick={editName}>Обновить</button>
                {userNameUpdated && <div>Имя успешно обновлено</div>}

                <br/><br/>
                Поменять аватарку:
                <br/>
                {!avatarUploading && <input type={"file"} onChange={onAvatarUpload}/>}
                {avatarUploading && <span>Загружаем аватар...</span>}
            </>}

            {user.posts?.items?.map(post => (
                <ComponentPost key={post.id} id={post.id}/>
            ))}

            {postsCursor && <button onClick={onShowMore}>Показать еще</button>}
        </div>
    )
}
