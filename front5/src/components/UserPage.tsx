import React, {useEffect} from "react";
import {api, User, UserPostsConnection} from "../store/types";
import {ComponentPost} from "./Post";
import styles from "./UserPage.module.css";
import produce from "immer";

export function UserPage() {
    const [user, setUser] = React.useState<User>();
    const [postsCursor, setPostsCursor] = React.useState("");
    const [viewerId, setViewerId] = React.useState("");
    const [loaded, setLoaded] = React.useState(false);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

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

    const localizeCounter = (count?: number) => {
        const mod = (count || 0) % 10;
        switch (mod) {
            case 0:
            case 5:
            case 6:
            case 7:
            case 8:
            case 9:
                return 'публикаций';
            case 1:
                return 'публикация';
            case 2:
            case 3:
            case 4:
                return 'публикации';
        }
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

    return (
        <div>
            <div className={styles.topBlock}>
                <img className={styles.userAvatar} src={user.avatar}/>
                <div className={styles.rightBlock}>
                    <div className={styles.userName}>{user.name}</div>
                    <div className={styles.userBio}>
                        {user.bio}
                    </div>
                    <div className={styles.userCounters}>
                        <div className={styles.userCounter}>
                            <b>{user.posts?.count}</b> {localizeCounter(user.posts?.count)}
                        </div>
                        {/*<div className={styles.userCounter}>
                            <b>9</b> подписчиков
                        </div>
                        <div className={styles.userCounter}>
                            <b>9</b> подписок
                        </div>*/}
                    </div>
                </div>
            </div>

            {user.id === viewerId && <>
                Имя: <input type="text" value={userName} onChange={e => setUserName(e.target.value)}/>
                <button onClick={editName}>Обновить</button>
                {userNameUpdated && <div>Имя успешно обновлено</div>}
            </>}

            {user.posts?.items?.map(post => (
                <ComponentPost key={post.id} post={post} onDelete={() => onPostDelete(post.id)}/>
            ))}

            {postsCursor && <button onClick={onShowMore}>Показать еще</button>}
        </div>
    )
}
