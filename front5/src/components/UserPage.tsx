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
        api("/users.list", {
            ids: location.pathname.substring(7),
            fields: "posts,posts.items,posts.items.user",
        }).then(r => {
            //const [user, viewerId] = r;
            setUser(r[0]);
            setViewerId('123');//TODO
            setUserName(r[0].name || "");

            setLoaded(true);
            setPostsCursor(r[0].posts?.nextCursor || "");
            setIsFollowing(Boolean(r[0].isFollowing))
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
        api("/users.list", {
            ids: location.pathname.substring(7),
            postsCursor: postsCursor,
            fields: "posts,posts.items,posts.items.user",
        }).then(r => {
            setPostsCursor(r[0].posts.nextCursor || "");

            setUser(produce(user, (user: User) => {
                user.posts = user.posts || {};
                user.posts.items = user.posts?.items || [];
                user.posts.items = [...user.posts.items, ...r[0].posts.items || []];
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
                <ComponentPost key={post.id} post={post} onDelete={() => onPostDelete(post.id)}/>
            ))}

            {postsCursor && <button onClick={onShowMore}>Показать еще</button>}
        </div>
    )
}
