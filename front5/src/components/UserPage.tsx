import React, {useEffect} from "react";
import {api, User} from "../store/types";
import {ComponentPost} from "./Post";
import styles from "./UserPage.module.css";

export function UserPage() {
    const [user, setUser] = React.useState<User>();
    const [viewerId, setViewerId] = React.useState("");
    const [loaded, setLoaded] = React.useState(false);

    const [userName, setUserName] = React.useState("");
    const [userNameUpdated, setUserNameUpdated] = React.useState(false);

    useEffect(() => {
        const f = new FormData();
        f.set("id", location.pathname.substring(7));

        api("/userPage", {
            id: location.pathname.substring(7)
        }).then(r => {
            const [user, viewerId] = r;
            setUser(user);
            setViewerId(viewerId);
            setUserName(user.name);

            // TODO strange
            for (let post of user.posts.posts) {
                post.user = user;
            }

            setLoaded(true);
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
            {user.posts?.posts.map(post => <ComponentPost key={post.id} post={post}/>)}
        </div>
    )
}
